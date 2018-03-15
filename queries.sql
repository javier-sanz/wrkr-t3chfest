SELECT date, sum(failure) as totalFailures from default.drive_stats
group by date, capacity_bytes

SELECT date, sum(failure) as totalFailures, (capacity_bytes /1024/1024/1024) GB  from default.drive_stats
group by date, capacity_bytes


SELECT date, GB, totalFailures/nDisk
FROM
  (SELECT date, sum(failure) AS totalFailures, (capacity_bytes /1024/1024/1024) GB,
    count(*) nDisk
   FROM default.drive_stats
   GROUP BY date, capacity_bytes)