CREATE (t:Town {name:"Казань"})

CREATE (pub:Publisher {
  id:1,
  name:"Полиграфический комбинат им. К. Якуба"
})
CREATE (pub)-[:LocatedAt]->(t)

CREATE (org:Organization {
  id:1,
  name:"Археологический кабинет Института истории им. Г. Ибрагимова КФАН СССР"
})
CREATE (org)-[:LocatedAt]->(t)

CREATE (a:Author {id:1, name:"Старостин П.Н."})
CREATE (r:Research {id:1, year:1985, description: "Памятники Западного Закамья Татарской АССР"})
CREATE (a)-[:Created]->(r)

CREATE (l:Literature {id:1, name:"Археологическая карта Татарской АССР - Западное Закамье"})
CREATE (r)-[:Has]->(l)
CREATE (org)-[:Stores {since: 1990, n:"435"}]->(l)
CREATE (pub)-[:Published {year: 1990}]->(l)

CREATE (:Epoch {id:1, name: 'Палеолит'})
CREATE (:Epoch {id:2, name: 'Мезолит'})
CREATE (:Epoch {id:3, name: 'Неолит'})
CREATE (:Epoch {id:4, name: 'Энеолит'})
CREATE (:Epoch {id:5, name: 'Бронзовый век'})
CREATE (:Epoch {id:6, name: 'Ранний железный век'})
CREATE (:Epoch {id:7, name: 'Великое переселение народов'})
CREATE (:Epoch {id:8, name: 'Средневековье'})
CREATE (:Epoch {id:9, name: 'Новое время'})