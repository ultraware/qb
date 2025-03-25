CREATE LOGIN qb_test WITH PASSWORD = 'qb_testA1', CHECK_POLICY = OFF;
CREATE USER qb_test FOR LOGIN qb_test;
EXEC sp_addrolemember 'db_owner', 'qb_test';
