// appears the most
MATCH (technology:Technology)<-[:TECHNOLOGY]-(reco)
RETURN technology.name, COUNT(*) AS appearances
ORDER BY appearances DESC;

// hold -> assess
MATCH (pos1:Position {value:"Hold"})<-[:POSITION]-(reco)-[:TECHNOLOGY]->(tech),
      (pos2:Position {value:"Assess"})<-[:POSITION]-(otherReco)-[:TECHNOLOGY]->(tech),
      (reco)-[:ON_DATE]->(recoDate),
      (otherReco)-[:ON_DATE]->(otherRecoDate)
WHERE (reco)-[:NEXT]->(otherReco)
RETURN tech.name AS technology, otherRecoDate.value AS dateOfChange;

// assess -> hold
MATCH (pos1:Position {value:"Assess"})<-[:POSITION]-(reco)-[:TECHNOLOGY]->(tech),
      (pos2:Position {value:"Hold"})<-[:POSITION]-(otherReco)-[:TECHNOLOGY]->(tech),
      (reco)-[:ON_DATE]->(recoDate),
      (otherReco)-[:ON_DATE]->(otherRecoDate)
WHERE (reco)-[:NEXT]->(otherReco)
RETURN tech.name AS technology, otherRecoDate.value AS dateOfChange;

// react
MATCH (t:Technology)<-[:TECHNOLOGY]-(reco)-[:ON_DATE]->(date), (reco)-[:POSITION]->(pos)
WHERE t.name contains "React.js"
RETURN pos.value, date.value
ORDER BY date.timestamp


// ember
MATCH (t:Technology)<-[:TECHNOLOGY]-(reco)-[:ON_DATE]->(date), (reco)-[:POSITION]->(pos)
WHERE t.name contains "Ember"
RETURN pos.value, date.value
ORDER BY date.timestamp;


// introduced in the latest radar
MATCH (date:Date {value: "Nov 2016"})<-[:ON_DATE]-(reco)-[:TECHNOLOGY]->(tech), (reco)-[:POSITION]->(position)
WHERE NOT (reco)<-[:NEXT]-()
WITH position, COUNT(*) AS count, COLLECT(tech.name) AS technologies
ORDER BY LENGTH((position)-[:NEXT*]->()) DESC
RETURN position.value, count, technologies;

