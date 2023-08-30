\c avito;
create table Users (user_id int PRIMARY key);
create table Segments (segment_name varchar(255) PRIMARY key);
create table SegmentsForUsers (
    user_id int references Users (user_id) on delete cascade,
    segment_name varchar(255) references Segments (segment_name) on delete cascade,
    PRIMARY KEY (user_id, segment_name)
);