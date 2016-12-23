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
