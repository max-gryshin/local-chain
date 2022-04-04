INSERT INTO public."user" (
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
    roles
) VALUES (
    'admin@mail.com',
    'max',
    'grishyn',
    'alexandrovich',
    '$2a$14$sy8dOPLdvgexL0U5Hvvpr.1Bds1VqxeY6TfM9RqCLYjS8B0uxQKSq', -- password 123123
    now(),
    now(),
    1,
    1,
    1,
    ARRAY['admin','manager']::role_type[]
);
