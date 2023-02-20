import { DataSourcePlugin } from '@grafana/data';
import { DataSource } from './datasource';
import { ConfigEditor } from './components/ConfigEditor';
import { QueryEditor } from './components/QueryEditor';
import { AmqpQuery, AmqpDataSourceOptions } from './types';

export const plugin = new DataSourcePlugin<DataSource, AmqpQuery, AmqpDataSourceOptions>(DataSource)
  .setConfigEditor(ConfigEditor)
  .setQueryEditor(QueryEditor);
