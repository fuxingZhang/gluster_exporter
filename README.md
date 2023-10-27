# gluster_exporter

## Usage

To run it:

```bash
./gluster_exporter [flags]
```

Help on flags:

```bash
./gluster_exporter -h
```

## Collectors

Collectors are enabled by providing a `--collector.<name>` flag.
Collectors that are enabled by default can be disabled by providing a `--no-collector.<name>` flag.

List:

| Collector      | Default | enable                     | disable                       |
| -------------- | ------- | -------------------------- | ----------------------------- |
| brick          | enable  | --collector.brick          | --no-collector.brick          |
| brick_status   | enable  | --collector.brick_status   | --no-collector.brick_status   |
| peer_counts    | enable  | --collector.peer_counts    | --no-collector.peer_counts    |
| peer_info      | enable  | --collector.peer_info      | --no-collector.peer_info      |
| ps             | enable  | --collector.ps             | --no-collector.ps             |
| volume_counts  | enable  | --collector.volume_counts  | --no-collector.volume_counts  |
| volume_heal    | enable  | --collector.volume_heal    | --no-collector.volume_heal    |
| volume_profile | enable  | --collector.volume_profile | --no-collector.volume_profile |
| volume_status  | enable  | --collector.volume_status  | --no-collector.volume_status  |

List of supported metrics are documented [here](./docs/metrics.md).
