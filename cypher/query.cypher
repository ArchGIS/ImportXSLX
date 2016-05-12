MATCH (au:Author {id:1}) // Автор исследования
MATCH (research:Research {id:1}) // Исследование, которое содержит карту
MATCH (map:Literature {id:1}) // Археологическая карта

// Эпохи
MATCH (medieval:Epoch {name:"Средневековье"}) 

// Памятник1: основные данные
CREATE (map)-[:References {pages:"54", n:"322"}]->(mon1)
CREATE (mon1:Monument {})
CREATE (medieval)-[:Has]->(mon1)
CREATE (k1:Knowledge{description:"...", name:"..."})
CREATE (research)-[:Contains]->(k1)
CREATE (k1)-[:Describes]->(mon1)

// Памятник1: Библиографическая ссылка1
MERGE (ra1:Author {name:"Федоров-Давыдов Г.А."})
MERGE (r1)-[:Created]->(r:Research {year: 1960})
CREATE (r1)-[:Has]->(l1:Literature {year: 1960})
CREATE (l1)-[:References {pages:"139", n:"55"}]->(mon1)

// Памятник1: Библиографическая ссылка2
MERGE (ra2:Author {name:"Фархутдинов Р.Г."})
MERGE (r2)-[:Created]->(r:Research {year: 1975})
CREATE (r2)-[:Has]->(l2:Literature {year: 1975})
CREATE (l2)-[:References {pages:"131", n:"818"}]->(mon1)
