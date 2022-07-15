import { Injectable, OnDestroy } from '@angular/core';
import { AgentGroup } from 'app/common/interfaces/orb/agent.group.interface';
import { Agent } from 'app/common/interfaces/orb/agent.interface';
import { AgentPolicy } from 'app/common/interfaces/orb/agent.policy.interface';
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
  merge,
  Observable,
  Subject,
  timer,
} from 'rxjs';
import {
  debounceTime,
  map,
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

  // race timer && forceRefresh until stopPolling
  private poller$: Observable<number>;

  pollController$: BehaviorSubject<boolean>;

  lastPollUpdate$: Subject<number>;

  // next to stop polling
  private killPolling: Subject<void>;

  // next to force refresh
  private forceRefresh: Subject<number>;

  // convenience polled observables
  // watch all pages available on agents
  private agents$: Observable<Agent[]>;
  private groups$: Observable<AgentGroup[]>;
  private datasets$: Observable<Dataset[]>;
  private policies$: Observable<AgentPolicy[]>;
  private sinks$: Observable<Sink[]>;

  private agentsTags$: Observable<string[]>;
  private groupsTags$: Observable<string[]>;
  private sinksTags$: Observable<string[]>;

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
    return this.poller$.pipe(
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

    this.poller$ = controller.pipe(takeUntil(this.killPolling));

    /**
     * TODO turn orb service into a poller service
     * available in root, inject it into a @Observe decorator
     * instead and wrap desired observables in it
     */

    this.agents$ = this.observe(this.agent.getAllAgents());

    this.groups$ = this.observe(
      this.group.getAllAgentGroups().pipe(map((page) => page.data)),
    );

    this.policies$ = this.observe(
      this.policy.getAllAgentPolicies().pipe(map((page) => page.data)),
    );

    this.datasets$ = this.observe(
      this.dataset.getAllDatasets().pipe(map((page) => page.data)),
    );

    this.sinks$ = this.observe(
      this.sink.getAllSinks().pipe(map((page) => page.data)),
    );

    this.agentsTags$ = this.agents$.pipe(
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

    const mapTags = (list: AgentGroup[] & Sink[]) => {
      return list
        .map((item) =>
          Object.entries(item.tags).map((entry) => `${entry[0]}: ${entry[1]}`),
        )
        .reduce((acc, val) => acc.concat(val), [])
        .filter(this.onlyUnique);
    };

    this.groupsTags$ = this.groups$.pipe(map((groups) => mapTags(groups)));

    this.sinksTags$ = this.sinks$.pipe(map((sinks) => mapTags(sinks)));
  }

  ngOnDestroy() {
    this.killPolling.next();
  }

  getAgentListView() {
    return this.agents$;
  }

  getAgentsTags() {
    return this.agentsTags$;
  }

  getGroupsTags() {
    return this.groupsTags$;
  }

  getGroupListView() {
    return this.groups$;
  }

  getPolicyListView() {
    return this.policies$;
  }

  getDatasetListView() {
    return this.datasets$;
  }

  getSinkListView() {
    return this.sinks$;
  }

  getSinksTags() {
    return this.sinksTags$;
  }

  onlyUnique = (value, index, self) => self.indexOf(value) === index;
}
