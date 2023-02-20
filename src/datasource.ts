import { DataSourceInstanceSettings } from '@grafana/data';
import { DataSourceWithBackend } from '@grafana/runtime';

import { AmqpQuery, AmqpDataSourceOptions } from './types';

export class DataSource extends DataSourceWithBackend<AmqpQuery, AmqpDataSourceOptions> {
  constructor(instanceSettings: DataSourceInstanceSettings<AmqpDataSourceOptions>) {
    super(instanceSettings);
  }
}
