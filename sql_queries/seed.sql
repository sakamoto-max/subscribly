CREATE OR REPLACE FUNCTION seed_users(user_count INT)
RETURNS VOID AS $$
BEGIN 
	INSERT INTO USERS(NAME, EMAIL, HASHED_PASSWORD, CREATED_AT, UPDATED_AT, JOINED_ORG)	
	SELECT 
		'User' || i,
		'user' || i || '@test.com',
		'$2a$107EqJtq98hPqEX7fNZaFWo05T0Kp7cH8FzZ3ZQZr8YF3P3p8F8e9cG',
		NOW(),
		NOW(),
		FALSE
	FROM generate_series(1, user_count) AS i;
END;
$$ LANGUAGE plpgsql;
-- --------------------------------------------
SELECT SEED_USES(500);
-------------------------------------------------
CREATE OR REPLACE FUNCTION seed_user_roles(
	start_user_id INT,
	end_user_id INT
)
RETURNS VOID AS $$
BEGIN
	INSERT INTO USER_ROLES(USER_ID, ROLE_ID)
	SELECT 
		ID, 
		3
	FROM USERS
	WHERE ID BETWEEN START_USER_ID AND END_USER_ID;
END;
$$ LANGUAGE plpgsql;
---------------------------------------------------
SELECT seed_user_roles(1, 500);
--------------------------------------------------
CREATE OR REPLACE FUNCTION SEED_ORGS_WITH_OWNERS(OWNER_IDS INT[])
RETURNs VOID AS $$
DECLARE
	OWNER_ID INT;
BEGIN 
	FOREACH OWNER_ID IN ARRAY OWNER_iDS
	LOOP
		INSERT INTO ORGS (ORG_NAME, OWNER_ID, CREATED_AT)
		VALUES(
			'org_user_' || owner_id,
			owner_id,
			NOW()
		);
	END LOOP;
END;
$$ LANGUAGE plpgsql;
-----------------------------------------------------
SELECT SEED_ORGS_WITH_OWNERS(ARRAY[50, 150, 200, 250, 300, 350, 400, 450, 500])
------------------------------------------------------
CREATE OR REPLACE FUNCTION SEED_ALL_USERS_INTO_ORGS(
	START_USER_ID INT,
	END_USER_ID INT,
	START_ORG_ID INT,
	END_ORG_ID INT
)
RETURNS VOID AS $$
DECLARE 
	ORG_COUNT INT := END_ORG_ID - START_ORG_ID + 1;
BEGIN 
	INSERT INTO ORG_MEMBERS (ORG_ID, USER_ID)
	SELECT 
		START_ORG_ID + ((U.ID - START_USER_ID) % ORG_COUNT) AS ORG_ID,
		U.ID
	FROM USERS U
	WHERE U.ID BETWEEN START_USER_ID AND END_USER_ID
		AND NOT EXISTS (
			SELECT 1
			FROM ORG_MEMBERS OM
			WHERE OM.USER_ID = U.ID
		);
END;
$$ LANGUAGE plpgsql;
--------------------------------------------------------
SELECT SEED_ALL_USERS_INTO_ORGS(27, 505, 10, 19)
--------------------------------------------------------
