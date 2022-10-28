CREATE TABLE IF NOT EXISTS users (
  id varchar(250) NOT NULL,
  nama varchar(250) NOT NULL,
  role varchar(250) NOT NULL,
  created_at INT NULL DEFAULT EXTRACT(epoch FROM NOW()),
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS harga (
  id varchar(250) NOT NULL,
  harga_topup INT NOT NULL,
  harga_buyback INT NOT NULL,
  created_by varchar(250) NOT NULL,
  created_at INT NULL DEFAULT EXTRACT(epoch FROM NOW()),
  PRIMARY KEY (id)
);


CREATE TABLE IF NOT EXISTS rekening (
  id varchar(250) NOT NULL,
  no_rek varchar(250) NOT NULL,
  saldo FLOAT NOT NULL,
  user_id varchar(250) NOT NULL,
  created_at INT NULL DEFAULT EXTRACT(epoch FROM NOW()),
  FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS top_up (
  id varchar(250) NOT NULL,
  harga INT NOT NULL,
  gram FLOAT NOT NULL,
  rekening_id varchar(250) NOT NULL,
  created_at INT NULL DEFAULT EXTRACT(epoch FROM NOW()),
  FOREIGN KEY (rekening_id) REFERENCES rekening (id) ON DELETE RESTRICT ON UPDATE CASCADE,
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS transaksi (
  id varchar(250) NOT NULL,
  type varchar(250) NOT NULL,
  top_up_id varchar(250) NULL,
  rekening_id varchar(250) NOT NULL,
  saldo FLOAT NOT NULL,
  gram FLOAT NOT NULL,
  harga_topup INT NOT NULL,
  harga_buyback INT NOT NULL,
  created_at INT NULL DEFAULT EXTRACT(epoch FROM NOW()),
  FOREIGN KEY (rekening_id) REFERENCES rekening (id) ON DELETE RESTRICT ON UPDATE CASCADE,
  PRIMARY KEY (id)
);

INSERT INTO users(id, nama, role, created_at) 
VALUES ('a001','admin1', 'admin', EXTRACT(epoch FROM NOW())),
('6pyXwTHVg','customer1', 'customer', EXTRACT(epoch FROM NOW()));

INSERT INTO rekening(id, no_rek, saldo, user_id, created_at) 
VALUES ('xzlcmANVR','898989', 0, '6pyXwTHVg',EXTRACT(epoch FROM NOW()));