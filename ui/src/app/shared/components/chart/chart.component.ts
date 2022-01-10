import { Component, Input, OnChanges, ViewChild } from '@angular/core';

import { ChartDataSets, ChartType, ChartOptions, ChartPoint } from 'chart.js';
import { BaseChartDirective, Color } from 'ng2-charts';
import { COLORS } from './chart.colors';
import { Dataset } from 'app/common/interfaces/mainflux.interface';
import { ToMillisecsPipe } from 'app/shared/pipes/time.pipe';

@Component({
  selector: 'ngx-chart',
  templateUrl: './chart.component.html',
  styleUrls: ['./chart.component.scss'],
})
export class ChartComponent implements OnChanges {
  chartColors: Color[] = COLORS;
  chartOptions: ChartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    elements: {
      line: {
        tension: 0.5,
      },
      point: {
        radius: 3,
      },
    },
    scales: {
      xAxes: [{
        type: 'time',
        distribution: 'series',
        ticks: {
          fontSize: 12,
          minRotation: 30,
        },
      }],
    },
  };

  chartDataSets: ChartDataSets[] = [];
  chartType: ChartType = 'scatter';

  @Input() msgDatasets: Dataset[] = [];
  @ViewChild(BaseChartDirective, { static: false }) chart: BaseChartDirective;
  constructor(
    private toMillisecsPipe: ToMillisecsPipe,
  ) { }

  ngOnChanges() {
    this.chartDataSets = [];

    this.msgDatasets.forEach( dataset => {
      const dataSet: ChartDataSets = {
        data: [],
        showLine: true,
        label: dataset.label,
      };

      // Create charts by name
      dataset.messages.forEach( msg => {
        const point: ChartPoint = {
          // Convert from seconds to milliseconds
          x: this.toMillisecsPipe.transform(msg.time),
          y: msg.value,
        };
        (dataSet.data as ChartPoint[]).push(point);
      });

      this.chartDataSets.push(dataSet);
      this.chart && this.chart.update();
    });
  }
}
