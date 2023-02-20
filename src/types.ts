import { DataQuery, DataSourceJsonData } from '@grafana/data';
//import internal from 'stream';

export interface AmqpQuery extends DataQuery {
  /*steamName: string;
  maxLengtrhBytes?: number;
  maxAge?: number;
  consumerName: string;
  crc?: boolean;
  constant: number;*/
  queryText?: string;
  constant: number;
}

export interface AmqpDataSourceOptions extends DataSourceJsonData {
  host: string;
  port: number;
  username?: string;
  tlsConnection?: boolean;
}

export interface AmqpSecureJsonData {
  password?: string;
}
