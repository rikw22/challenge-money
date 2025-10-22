
-- DDL 
CREATE TABLE account (
  ID INTEGER PRIMARY KEY,
  document_number VARCHAR(11)
);

CREATE TABLE operationtype (
  ID INTEGER PRIMARY KEY,
  description VARCHAR(50)
);
 
CREATE TABLE transaction (
  ID UUID PRIMARY KEY DEFAULT uuidv7(),
  account_id INTEGER REFERENCES account(ID),
  operationtype_id INTEGER REFERENCES operationtype(ID),
  amount INTEGER,
  eventdate TIMESTAMP
);


-- DML - Initial data
-- Account
INSERT INTO account VALUES(1, '11111111111'), (2, '22222222222');

-- Operation Types
INSERT INTO operationtype VALUES(1, 'Normal Purchase'), (2, 'Purchase with installments'), (3, 'Withdrawal'), (4, 'Credit Voucher');

-- Transaction
INSERT INTO transaction (ID, account_id, operationtype_id, amount, eventdate) VALUES
('019a096b-ad9f-7f0e-88a4-9c93a754b029', 1, 1, -5000, '2020-01-01T10:32:07.7199222'),
('019a096e-d76f-75c0-8ae1-cca9a4b99cd7', 1, 4, 6000, '2020-01-01T10:32:07.7199222');

