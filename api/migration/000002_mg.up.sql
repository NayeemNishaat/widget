-- SQLINES DEMO *** 9  Distrib 10.5.9-MariaDB, for osx10.16 (x86_64)
--
-- SQLINES DEMO ***   Database: widgets
-- SQLINES DEMO *** -------------------------------------
-- SQLINES DEMO *** .5.9-MariaDB

/* SQLINES DEMO *** ARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/* SQLINES DEMO *** ARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/* SQLINES DEMO *** LLATION_CONNECTION=@@COLLATION_CONNECTION */;
/* SQLINES DEMO *** tf8mb4 */;
/* SQLINES DEMO *** ME_ZONE=@@TIME_ZONE */;
/* SQLINES DEMO *** NE='+00:00' */;
/* SQLINES DEMO *** IQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/* SQLINES DEMO *** REIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/* SQLINES DEMO *** L_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/* SQLINES DEMO *** L_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- SQLINES DEMO *** or table `customers`
--
SET session_replication_role = 'replica';

DROP TABLE IF EXISTS customers;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE customers (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
)  ;

ALTER SEQUENCE customers_id_seq RESTART WITH 4;
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;

--
-- SQLINES DEMO *** or table `statuses`
--

DROP TABLE IF EXISTS statuses;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE statuses (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  name varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
)  ;

ALTER SEQUENCE statuses_id_seq RESTART WITH 4;
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;

--
-- SQLINES DEMO *** or table `transaction_statuses`
--

DROP TABLE IF EXISTS transaction_statuses;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE transaction_statuses (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  name varchar(255) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
)  ;

ALTER SEQUENCE transaction_statuses_id_seq RESTART WITH 6;
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;

--
-- SQLINES DEMO *** or table `transactions`
--

DROP TABLE IF EXISTS transactions;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE transactions (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  amount int NOT NULL,
  currency varchar(255) NOT NULL,
  last_four varchar(255) NOT NULL,
  bank_return_code varchar(255) NOT NULL,
  transaction_status_id int NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  expiry_month int NOT NULL DEFAULT 0,
  expiry_year int NOT NULL DEFAULT 0,
  payment_intent varchar(255) NOT NULL DEFAULT '',
  payment_method varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (id)
,
  CONSTRAINT transactions_transaction_statuses_id_fk FOREIGN KEY (transaction_status_id) REFERENCES transaction_statuses (id) ON DELETE CASCADE ON UPDATE CASCADE
)  ;

ALTER SEQUENCE transactions_id_seq RESTART WITH 4;
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;

CREATE INDEX transactions_transaction_statuses_id_fk ON transactions (transaction_status_id);

--
-- SQLINES DEMO *** or table `users`
--

DROP TABLE IF EXISTS users;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE users (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  first_name varchar(255) NOT NULL,
  last_name varchar(255) NOT NULL,
  email varchar(255) NOT NULL,
  password varchar(60) NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
)  ;

ALTER SEQUENCE users_id_seq RESTART WITH 2;
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;

--
-- SQLINES DEMO *** or table `widgets`
--

DROP TABLE IF EXISTS widgets;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE widgets (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  name varchar(255) NOT NULL DEFAULT '',
  description text NOT NULL DEFAULT '',
  inventory_level int NOT NULL,
  price int NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  image varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (id)
)  ;

ALTER SEQUENCE widgets_id_seq RESTART WITH 2;

--
-- SQLINES DEMO *** or table `orders`
--

DROP TABLE IF EXISTS orders;
/* SQLINES DEMO *** cs_client     = @@character_set_client */;
/* SQLINES DEMO *** er_set_client = utf8 */;
-- SQLINES LICENSE FOR EVALUATION USE ONLY
CREATE TABLE orders (
  id int NOT NULL GENERATED ALWAYS AS IDENTITY,
  widget_id int NOT NULL,
  transaction_id int NOT NULL,
  status_id int NOT NULL,
  quantity int NOT NULL,
  amount int NOT NULL,
  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  customers_id int NOT NULL,
  PRIMARY KEY (id)
,
  CONSTRAINT orders_customers_id_fk FOREIGN KEY (customers_id) REFERENCES customers (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT orders_statuses_id_fk FOREIGN KEY (status_id) REFERENCES statuses (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT orders_transactions_id_fk FOREIGN KEY (transaction_id) REFERENCES transactions (id) ON DELETE CASCADE ON UPDATE CASCADE,
  CONSTRAINT orders_widgets_id_fk FOREIGN KEY (widget_id) REFERENCES widgets (id) ON DELETE CASCADE ON UPDATE CASCADE
)  ;

ALTER SEQUENCE orders_id_seq RESTART WITH 2;
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;

CREATE INDEX orders_widgets_id_fk ON orders (widget_id);
CREATE INDEX orders_transactions_id_fk ON orders (transaction_id);
CREATE INDEX orders_statuses_id_fk ON orders (status_id);
CREATE INDEX orders_customers_id_fk ON orders (customers_id);
/* SQLINES DEMO *** er_set_client = @saved_cs_client */;
/* SQLINES DEMO *** NE=@OLD_TIME_ZONE */;

/* SQLINES DEMO *** E=@OLD_SQL_MODE */;
/* SQLINES DEMO *** _KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/* SQLINES DEMO *** CHECKS=@OLD_UNIQUE_CHECKS */;
/* SQLINES DEMO *** ER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/* SQLINES DEMO *** ER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/* SQLINES DEMO *** ON_CONNECTION=@OLD_COLLATION_CONNECTION */;
/* SQLINES DEMO *** ES=@OLD_SQL_NOTES */;

-- SQLINES DEMO ***  2021-07-15 14:24:57

SET session_replication_role = 'origin';