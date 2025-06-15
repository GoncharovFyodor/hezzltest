CREATE TABLE goods
(
    Id UInt32,
    ProjectId UInt32,
    Name String,
    Description String,
    Priority UInt32,
    Removed UInt8,
    EventTime DateTime
) ENGINE = MergeTree()
ORDER BY (Id, ProjectId, Name);