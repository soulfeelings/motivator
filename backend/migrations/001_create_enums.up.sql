CREATE TYPE user_role AS ENUM ('owner', 'admin', 'manager', 'employee');
CREATE TYPE invite_status AS ENUM ('pending', 'accepted', 'expired', 'revoked');
