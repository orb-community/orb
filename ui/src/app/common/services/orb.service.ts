import { Injectable, OnDestroy } from '@angular/core';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Dataset } from 'app/common/interfaces/orb/dataset.policy.interface';
import { Sink } from 'app/common/interfaces/orb/sink.interface';
import { AgentGroupsService } from 'app/common/services/agents/agent.groups.service';
import { AgentPoliciesService } from 'app/common/services/agents/agent.policies.service';
import { AgentsService } from 'app/common/services/agents/agents.service';
import { DatasetPoliciesService } from 'app/common/services/dataset/dataset.policies.service';
import { SinksService } from 'app/common/services/sinks/sinks.service';
import {
  BehaviorSubject,
  defer,
  EMPTY,
  forkJoin,
  merge,
  Observable,
  Subject,
  timer,
} from 'rxjs';
import {
  debounceTime,
  map,
  mergeMap,
  retry,
  shareReplay,
  switchMap,
  takeUntil,
  tap,
} from 'rxjs/operators';

export const PollControls = {
  PAUSE: false,
  RESUME: true,
};

@Injectable({
  providedIn: 'root',
})
export class OrbService implements OnDestroy {
  // interval for timer
  pollInterval = 1000;

  pollController$: BehaviorSubject<boolean>;

  lastPollUpdate$: Subject<number>;

  // next to stop polling
  killPolling: Subject<void>;

  // next to force refresh
  private forceRefresh: Subject<number>;

  pausePolling() {
    this.pollController$.next(PollControls.PAUSE);
  }

  startPolling() {
    this.pollController$.next(PollControls.RESUME);
  }

  refreshNow() {
    this.forceRefresh.next(1);
  }

  observe<T>(observable: Observable<T>) {
    const controller = merge(
      this.pollController$.pipe(
        switchMap((control) => {
          if (control === PollControls.RESUME)
            return defer(() => timer(1, this.pollInterval));
          return EMPTY;
        }),
      ),
      this.forceRefresh.pipe(debounceTime(1000)),
    );

    const poller$ = controller.pipe(takeUntil(this.killPolling));

    return poller$.pipe(
      switchMap(() =>
        observable.pipe(
          tap((_) => {
            this.lastPollUpdate$.next(Date.now());
          }),
        ),
      ),
      retry(),
      shareReplay(1),
    );
  }

  constructor(
    private agent: AgentsService,
    private dataset: DatasetPoliciesService,
    private group: AgentGroupsService,
    private policy: AgentPoliciesService,
    private sink: SinksService,
  ) {
    this.lastPollUpdate$ = new Subject<number>();
    this.forceRefresh = new Subject<number>();
    this.killPolling = new Subject<void>();

    this.pollController$ = new BehaviorSubject<boolean>(PollControls.PAUSE);
  }

  private mapTags = (list: AgentGroup[] & Sink[]) => {
    return list
      .map((item) =>
        Object.entries(item.tags).map((entry) => `${entry[0]}: ${entry[1]}`),
      )
      .reduce((acc, val) => acc.concat(val), [])
      .filter(this.onlyUnique);
  }

  ngOnDestroy() {
    this.killPolling.next();
  }

  getAgentListView() {
    return this.observe(this.agent.getAllAgents());
  }

  getAgentsTags() {
    return this.observe(this.agent.getAllAgents()).pipe(
      map((agents) =>
        agents
          .map((_agent) =>
            Object.entries(_agent.orb_tags)
              .map((entry) => `${entry[0]}: ${entry[1]}`)
              .concat(
                Object.entries(_agent.agent_tags).map(
                  (entry) => `${entry[0]}: ${entry[1]}`,
                ),
              ),
          )
          .reduce((acc, val) => acc.concat(val), [])
          .filter(this.onlyUnique),
      ),
    );
  }

  getGroupsTags() {
    return this.observe(this.group.getAllAgentGroups()).pipe(
      map((groups) => this.mapTags(groups)),
    );
  }

  getGroupListView() {
    return this.observe(this.group.getAllAgentGroups());
  }

  getPolicyListView() {
    return this.observe(this.policy.getAllAgentPolicies());
  }

  getPolicyFullView(id: string) {
    // retrieve policy
    return this.policy.getAgentPolicyById(id).pipe(
      mergeMap((policy) =>
        // need a way to get a dataset linked to a policy without having to filter it out
        this.dataset.getAllDatasets().pipe(
          map((_dataset) =>
            _dataset.filter((dataset) => policy.id === dataset.agent_policy_id),
          ),
          // from the filtered dataset list, query all agent groups associated with the list
          mergeMap((datasets: Dataset[]) =>
            forkJoin(
              datasets
                .map((dataset) => dataset?.agent_group_id)
                .map((groupId) => this.group.getAgentGroupById(groupId)),
            ).pipe(map((groups) => ({ datasets, groups, policy }))),
          ),
          // same for sinks
          mergeMap(({ datasets, groups }) =>
            forkJoin(
              datasets
                .map((dataset) => dataset?.sink_ids)
                .reduce((acc, val) => acc.concat(val), [])
                .map((sinkId) => this.sink.getSinkById(sinkId)),
            ).pipe(map((sinks) => ({ datasets, sinks, policy, groups }))),
          ),
        ),
      ),
      // from here on I can map to any shape I like
      // dataset list uses the info below
      map(({ datasets, sinks, policy, groups }) => ({
        datasets: datasets.map((dataset) => ({
          ...dataset,
          agent_group: groups.find(
            (group) => group.id === dataset.agent_group_id,
          ),
          agent_policy: policy,
          sinks: sinks.filter((sink) => dataset.sink_ids.includes(sink.id)),
        })),
        sinks,
        policy: { ...policy, groups, datasets },
        groups,
      })),
    );
  }

  getDatasetListView() {
    return this.observe(this.dataset.getAllDatasets());
  }

  getSinkListView() {
    return this.observe(this.sink.getAllSinks());
  }

  getSinksTags() {
    return this.observe(this.sink.getAllSinks()).pipe(
      map((sinks) => this.mapTags(sinks)),
    );
  }

  onlyUnique = (value, index, self) => self.indexOf(value) === index;
}
