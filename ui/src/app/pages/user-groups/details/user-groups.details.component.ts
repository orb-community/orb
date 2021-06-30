import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

import { UserGroup, User } from 'app/common/interfaces/mainflux.interface';
import { UsersService } from 'app/common/services/users/users.service';
import { UserGroupsService } from 'app/common/services/users/groups.service';
import { NotificationsService } from 'app/common/services/notifications/notifications.service';

@Component({
  selector: 'ngx-user-groups-details-component',
  templateUrl: './user-groups.details.component.html',
  styleUrls: ['./user-groups.details.component.scss'],
})
export class UserGroupsDetailsComponent implements OnInit {
  offset = 0;
  limit = 20;

  userGroup: UserGroup = {};
  users: User[] = [];
  members: User[] = [];

  selectedUsers = [];

  constructor(
    private route: ActivatedRoute,
    private usersService: UsersService,
    private userGroupsService: UserGroupsService,
    private notificationsService: NotificationsService,
  ) {}

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');

    this.userGroupsService.getGroup(id).subscribe(
      (resp: any) => {
        this.userGroup = resp;

        this.getMembers();
      },
    );
  }

  getMembers() {
    this.usersService.getUsers().subscribe(
      (resp: any) => {
        this.users = resp.users;
      },
    );

    this.userGroupsService.getMembers(this.userGroup.id).subscribe(
      respMemb => {
        this.members = respMemb.users;

        if (this.members) {
          // Remove members from available Users
          this.members.forEach(m => {
            this.users = this.users.filter(u => u.id !== m.id);
          });
        }
      },
    );
  }

  onAssign() {
    const userIDs = this.selectedUsers.map(u => u.id);

    this.userGroupsService.assignUser(this.userGroup.id, userIDs).subscribe(
      resp => {
        this.notificationsService.success('Successfully assigned User(s) to Group', '');
        this.selectedUsers = [];
        this.getMembers();
      },
    );

    if (this.selectedUsers.length === 0) {
      this.notificationsService.warn('User(s) must be provided', '');
    }
  }

  onUnassign(member: any) {
    this.userGroupsService.unassignUser(this.userGroup.id, [member.id]).subscribe(
      resp => {
        this.notificationsService.success('Successfully unassigned User(s) from Group', '');
        this.selectedUsers = [];
        this.getMembers();
      },
    );
  }
}
