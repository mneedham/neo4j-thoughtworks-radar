WITH ["Hold", "Assess", "Trial", "Adopt"] AS positions
UNWIND RANGE (0, size(positions) - 2) AS index
WITH positions[index] AS pos1, positions[index + 1] AS pos2
MERGE (position1:Position {value: pos1})
MERGE (position2:Position {value: pos2})
MERGE (position1)-[:NEXT]->(position2);


load csv with headers from "file:///blips.csv" AS row
MATCH (position:Position {value:  row.suggestion })
MERGE (tech:Technology {name:  row.technology })
MERGE (date:Date {value: row.date})
MERGE (recommendation:Recommendation {id: tech.name + "_" + date.value + "_" + position.value})
MERGE (recommendation)-[:ON_DATE]->(date)
MERGE (recommendation)-[:POSITION]->(position)
MERGE (recommendation)-[:TECHNOLOGY]->(tech);

match (date:Date)
SET date.timestamp = apoc.date.parse(date.value, "ms", "MMM yyyy");

MATCH (date:Date)
WITH date
ORDER BY date.timestamp
WITH COLLECT(date) AS dates
UNWIND range(0, size(dates)-2) AS index
WITH dates[index] as month1, dates[index+1] AS month2
MERGE (month1)-[:NEXT]->(month2);

MATCH (tech)<-[:TECHNOLOGY]-(reco:Recommendation)-[:ON_DATE]->(date)
WITH tech, reco, date
ORDER BY tech.name, date.timestamp
WITH tech, COLLECT(reco) AS recos
UNWIND range(0, size(recos)-2) AS index
WITH recos[index] AS reco1, recos[index+1] AS reco2
MERGE (reco1)-[:NEXT]->(reco2);
