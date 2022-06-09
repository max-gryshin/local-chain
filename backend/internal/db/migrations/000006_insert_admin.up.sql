INSERT INTO public."user" (
    id,
    email,
    first_name,
    last_name,
    middle_name,
    password_hash,
    created_at,
    updated_at,
    created_by,
    updated_by,
    status,
    roles,
    manager_id
) VALUES (
    1,
    'admin@mail.com',
    'admin',
    'admin',
    'admin',
    '$2a$14$sy8dOPLdvgexL0U5Hvvpr.1Bds1VqxeY6TfM9RqCLYjS8B0uxQKSq', -- password 123123
    now(),
    now(),
    1,
    1,
    1,
    ARRAY['admin','manager']::role_type[],
    1
);
