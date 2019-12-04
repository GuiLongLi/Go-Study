DROP TABLE IF EXISTS user;
CREATE TABLE user (
  id bigint(20) NOT NULL AUTO_INCREMENT,
  user_name varchar(20) DEFAULT "",
  password varchar(20) DEFAULT "",
  age tinyint(11) DEFAULT "0",
  PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4;
INSERT INTO user VALUES ("1", "admin", "12164","18");