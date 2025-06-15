CREATE TABLE IF NOT EXISTS default.goods
(
    Id UInt32,
    ProjectId UInt32,
    Name String,
    Description String,
    Priority UInt32,
    Removed Bool,
    EventTime DateTime
) ENGINE = MergeTree()
ORDER BY (Id, ProjectId, Name);