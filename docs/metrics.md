# Metrics Exported by Gluster Prometheus exporter

## collector.brick

### gluster_brick_capacity_used_bytes

Used capacity of gluster bricks in bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |

### gluster_brick_capacity_free_bytes

Free capacity of gluster bricks in bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |

### gluster_brick_capacity_bytes_total

Total capacity of gluster bricks in bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |

### gluster_brick_inodes_total

Total no of inodes of gluster brick disk

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |

### gluster_brick_inodes_free

Free no of inodes of gluster brick disk

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |

### gluster_brick_inodes_used

Used no of inodes of gluster brick disk

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |

### gluster_subvol_capacity_used_bytes

Effective used capacity of gluster subvolume in bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| volume     | Volume Name     |
| subvolume  | Sub volume name |

### gluster_subvol_capacity_total_bytes

Effective total capacity of gluster subvolume in bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| volume     | Volume Name     |
| subvolume  | Sub volume name |

### gluster_brick_lv_size_bytes

Bricks LV size Bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |
| vg_name    | VG Name         |
| lv_path    | LV Path         |
| lv_uuid    | UUID of LV      |

### gluster_brick_lv_percent

Bricks LV usage percent

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |
| vg_name    | VG Name         |
| lv_path    | LV Path         |
| lv_uuid    | UUID of LV      |

### gluster_brick_lv_metadata_size_bytes

Bricks LV metadata size Bytes

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |
| vg_name    | VG Name         |
| lv_path    | LV Path         |
| lv_uuid    | UUID of LV      |

### gluster_brick_lv_metadata_percent

Bricks LV metadata usage percent

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |
| vg_name    | VG Name         |
| lv_path    | LV Path         |
| lv_uuid    | UUID of LV      |

### gluster_vg_extent_total_count

VG extent total count

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |
| vg_name    | VG Name         |
| lv_path    | LV Path         |
| lv_uuid    | UUID of LV      |

### gluster_vg_extent_alloc_count

VG extent allocated count

| Label      | Description     |
| ---------- | --------------- |
| cluster_id | Cluster ID      |
| host       | Host name or IP |
| id         | Brick ID        |
| brick_path | Brick Path      |
| volume     | Volume Name     |
| subvolume  | Sub Volume name |
| vg_name    | VG Name         |
| lv_path    | LV Path         |
| lv_uuid    | UUID of LV      |

### gluster_thinpool_data_total_bytes

Thin pool size Bytes

| Label         | Description              |
| ------------- | ------------------------ |
| cluster_id    | Cluster ID               |
| host          | Host name or IP          |
| thinpool_name | Name of the thinpool LV  |
| vg_name       | Name of the Volume Group |
| volume        | Volume Name              |
| subvolume     | Name of the Subvolume    |
| brick_path    | Brick Path               |

### gluster_thinpool_data_used_bytes

Thin pool data used Bytes

| Label         | Description              |
| ------------- | ------------------------ |
| cluster_id    | Cluster ID               |
| host          | Host name or IP          |
| thinpool_name | Name of the thinpool LV  |
| vg_name       | Name of the Volume Group |
| volume        | Volume Name              |
| subvolume     | Name of the Subvolume    |
| brick_path    | Brick Path               |

### gluster_thinpool_metadata_total_bytes

Thin pool metadata size Bytes

| Label         | Description              |
| ------------- | ------------------------ |
| cluster_id    | Cluster ID               |
| host          | Host name or IP          |
| thinpool_name | Name of the thinpool LV  |
| vg_name       | Name of the Volume Group |
| volume        | Volume Name              |
| subvolume     | Name of the Subvolume    |
| brick_path    | Brick Path               |

### gluster_thinpool_metadata_used_bytes

Thin pool metadata used Bytes

| Label         | Description              |
| ------------- | ------------------------ |
| cluster_id    | Cluster ID               |
| host          | Host name or IP          |
| thinpool_name | Name of the thinpool LV  |
| vg_name       | Name of the Volume Group |
| volume        | Volume Name              |
| subvolume     | Name of the Subvolume    |
| brick_path    | Brick Path               |

## collector.brick_status

### gluster_brick_up

Brick up (1-up, 0-down)

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume Name         |
| hostname   | Host name or IP     |
| brick_path | Brick Path          |
| peer_id    | Peer ID             |
| pid        | Process ID of brick |

## collector.peer_counts

### gluster_pv_count

No: of Physical Volumes

| Label      | Description                                           |
| ---------- | ----------------------------------------------------- |
| cluster_id | Cluster ID                                            |
| name       | Metric name, for which data is collected              |
| peer_id    | Peer ID of the host on which this metric is collected |

### gluster_lv_count

No: of Logical Volumes in a Volume Group

| Label      | Description                                           |
| ---------- | ----------------------------------------------------- |
| cluster_id | Cluster ID                                            |
| name       | Metric name, for which the data is collected          |
| peer_id    | Peer ID of the host on which this metric is collected |
| vg_name    | Volume Group Name associated with the metric          |

### gluster_vg_count

No: of Volume Groups

| Label      | Description                                           |
| ---------- | ----------------------------------------------------- |
| cluster_id | Cluster ID                                            |
| name       | Metric name, for which data is collected              |
| peer_id    | Peer ID of the host on which this metric is collected |

### gluster_thinpool_count

No: of thinpools in a Volume Group

| Label      | Description                                           |
| ---------- | ----------------------------------------------------- |
| cluster_id | Cluster ID                                            |
| name       | Metric name, for which the data is collected          |
| peer_id    | Peer ID of the host on which this metric is collected |
| vg_name    | Volume Group Name associated with the metric          |

## collector.peer_info

### gluster_peer_count

Number of peers in cluster

| Label    | Description                                                       |
| -------- | ----------------------------------------------------------------- |
| instance | Hostname of the gluster-prometheus instance providing this metric |

### gluster_peer_status

Peer status info

| Label    | Description                                                       |
| -------- | ----------------------------------------------------------------- |
| instance | Hostname of the gluster-prometheus instance providing this metric |
| hostname | Hostname of the peer for which data is collected                  |
| uuid     | Uuid of the peer for which data is collected                      |

### gluster_peer_connected

Peer connection status

| Label    | Description                                                       |
| -------- | ----------------------------------------------------------------- |
| instance | Hostname of the gluster-prometheus instance providing this metric |
| hostname | Hostname of the peer for which data is collected                  |
| uuid     | Uuid of the peer for which data is collected                      |

## collector.ps

### gluster_cpu_percentage

CPU percentage of Gluster process. One metric will be exposed for each process. Note: values of labels will be empty if not applicable to that process. For example, glusterd process will not have labels for volume or brick_path. It is the CPU time used divided by the time the process has been running (cputime/realtime ratio), expressed as a percentage.

| Label      | Description                                                   |
| ---------- | ------------------------------------------------------------- |
| cluster_id | Cluster ID                                                    |
| volume     | Volume Name                                                   |
| peer_id    | Peer ID                                                       |
| brick_path | Brick Path                                                    |
| name       | Name of the Gluster process(Ex: `glusterfsd`, `glusterd` etc) |

### gluster_memory_percentage

Memory percentage of Gluster process. One metric will be exposed for each process. Note: values of labels will be empty if not applicable to that process. For example, glusterd process will not have labels for volume or brick_path. It is the ratio of the process’s resident set size to the physical memory on the machine, expressed as a percentage

| Label      | Description                                                   |
| ---------- | ------------------------------------------------------------- |
| cluster_id | Cluster ID                                                    |
| volume     | Volume Name                                                   |
| peer_id    | Peer ID                                                       |
| brick_path | Brick Path                                                    |
| name       | Name of the Gluster process(Ex: `glusterfsd`, `glusterd` etc) |

### gluster_resident_memory_bytes

Resident Memory of Gluster process in bytes. One metric will be exposed for each process. Note: values of labels will be empty if not applicable to that process. For example, glusterd process will not have labels for volume or brick_path.

| Label      | Description                                                   |
| ---------- | ------------------------------------------------------------- |
| cluster_id | Cluster ID                                                    |
| volume     | Volume Name                                                   |
| peer_id    | Peer ID                                                       |
| brick_path | Brick Path                                                    |
| name       | Name of the Gluster process(Ex: `glusterfsd`, `glusterd` etc) |

### gluster_virtual_memory_bytes

Virtual Memory of Gluster process in bytes. One metric will be exposed for each process. Note: values of labels will be empty if not applicable to that process. For example, glusterd process will not have labels for volume or brick_path.

| Label      | Description                                                   |
| ---------- | ------------------------------------------------------------- |
| cluster_id | Cluster ID                                                    |
| volume     | Volume Name                                                   |
| peer_id    | Peer ID                                                       |
| brick_path | Brick Path                                                    |
| name       | Name of the Gluster process(Ex: `glusterfsd`, `glusterd` etc) |

### gluster_elapsed_time_seconds

Elapsed Time or Uptime of Gluster processes in seconds. One metric will be exposed for each process. Note: values of labels will be empty if not applicable to that process. For example, glusterd process will not have labels for volume or brick_path.

| Label      | Description                                                   |
| ---------- | ------------------------------------------------------------- |
| cluster_id | Cluster ID                                                    |
| volume     | Volume Name                                                   |
| peer_id    | Peer ID                                                       |
| brick_path | Brick Path                                                    |
| name       | Name of the Gluster process(Ex: `glusterfsd`, `glusterd` etc) |

## collector.volume_counts

### gluster_volume_total_count

Total no of volumes

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |

### gluster_volume_created_count

Freshly created no of volumes

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |

### gluster_volume_started_count

Total no of started volumes

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |

### gluster_volume_brick_count

Total no of bricks in volume

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume Name |

### gluster_volume_snapshot_brick_count_total

Total count of snapshots bricks for volume

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume Name |

### gluster_volume_snapshot_brick_count_active

Total active count of snapshots bricks for volume

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume Name |

### gluster_volume_up

Volume is started or not (1-started, 0-not started)

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume Name |

## collector.volume_heal

### gluster_volume_heal_count

self heal count for volume

| Label      | Description    |
| ---------- | -------------- |
| cluster_id | Cluster ID     |
| volume     | Volume Name    |
| brick_path | Brick Path     |
| host       | Hostname or IP |

### gluster_volume_split_brain_heal_count

self heal count for volume in split brain

| Label      | Description    |
| ---------- | -------------- |
| cluster_id | Cluster ID     |
| volume     | Volume Name    |
| brick_path | Brick Path     |
| host       | Hostname or IP |

## collector.volume_profile

### gluster_volume_profile_total_reads

Total no of reads

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume name |
| brick      | Brick Name  |

### gluster_volume_profile_total_writes

Total no of writes

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume name |
| brick      | Brick Name  |

### gluster_volume_profile_duration_secs

Duration

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume name |
| brick      | Brick Name  |

### gluster_volume_profile_total_reads_interval

Total no of reads for interval stats

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume name |
| brick      | Brick Name  |

### gluster_volume_profile_total_writes_interval

Total no of writes for interval stats

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume name |
| brick      | Brick Name  |

### gluster_volume_profile_duration_secs_interval

Duration for interval stats

| Label      | Description |
| ---------- | ----------- |
| cluster_id | Cluster ID  |
| volume     | Volume name |
| brick      | Brick Name  |

### gluster_volume_profile_fop_hits

Cumulative FOP hits

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_avg_latency

Cumulative FOP avergae latency

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_min_latency

Cumulative FOP min latency

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_max_latency

Cumulative FOP max latency

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_hits_interval

Interval based FOP hits

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_avg_latency_interval

Interval based FOP average latency

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_min_latency_interval

Interval based FOP min latency

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_max_latency_interval

Interval based FOP max latency

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_total_hits_on_aggregated_fops

Cumulative total hits on aggregated FOPs like READ_WRIET_OPS, LOCK_OPS, INODE_OPS etc

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

### gluster_volume_profile_fop_total_hits_on_aggregated_fops_interval

Interval based total hits on aggregated FOPs like READ_WRIET_OPS, LOCK_OPS, INODE_OPS etc

| Label      | Description         |
| ---------- | ------------------- |
| cluster_id | Cluster ID          |
| volume     | Volume name         |
| brick      | Brick Name          |
| host       | Hostname or IP      |
| fop        | File Operation name |

## collector.volume_status

### gluster_volume_status_brick_count

Number of bricks for volume

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |

### gluster_volume_brick_status

Per node brick status for volume

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |

### gluster_volume_brick_port

Brick port

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |

### gluster_volume_brick_pid

Brick pid

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |

### gluster_volume_brick_total_inodes

Brick total inodes

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |

### gluster_volume_brick_free_inodes

Brick free inodes

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |

### gluster_volume_brick_total_bytes

Brick total bytes

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |

### gluster_volume_brick_free_bytes

Brick free bytes

| Label       | Description                                                       |
| ----------- | ----------------------------------------------------------------- |
| instance    | Hostname of the gluster-prometheus instance providing this metric |
| volume_name | Name of the volume                                                |
| hostname    | Hostname of the brick                                             |
| peer_id     | Uuid of the peer hosting this brick                               |
