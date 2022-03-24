alter table "user"
    add constraint user_created_by
        foreign key (created_by) references "user",
    add constraint user_updated_by
        foreign key (updated_by) references "user";

alter table account
    add constraint account_user_id_fk
        foreign key (user_id) references "user",
    add constraint account_parent
        foreign key (manager_id) references account,
    add constraint account_user_created_by
        foreign key (created_by) references "user",
    add constraint account_user_updated_by
        foreign key (updated_by) references "user";

alter table wallet
    add constraint wallet_account_id_fk
        foreign key (account_id) references account,
    add constraint wallet_user_created_by
        foreign key (created_by) references "user",
    add constraint wallet_user_updated_by
        foreign key (updated_by) references "user";

alter table "order"
    add constraint order_wallet_id_fk
        foreign key (wallet_id) references wallet (id),
    add constraint order_user_created_by
        foreign key (created_by) references "user",
    add constraint order_user_updated_by
        foreign key (updated_by) references "user";