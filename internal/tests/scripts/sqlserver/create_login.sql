CREATE LOGIN qb_test WITH PASSWORD = 'qb_testA1';
CREATE USER qb_test FOR LOGIN qb_test;
EXEC sp_addrolemember 'db_owner', 'qb_test';
