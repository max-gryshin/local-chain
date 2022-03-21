API
    user:
        1. user list.
            Request:
                `GET /api/user/all?page=1&limit=30`
            Response:
            `{
                {
                    "name": "max",
                    "dob": 1993.11.11,
                    "status": "active"                
                },
                {
                    "name": "Jhon",
                    "dob": 1985.05.04,
                    "status": "active"                
                },
                ...
            }`
        2. certain user
            Request:
                `GET /api/user/{id}`
            Response:
                `{
                    "name": "max",
                    "dob": 1993.11.11,
                    "status": "active"
                    ...
                }`
    account:
        1. user account
            Request:
                `GET /api/user/{id}/account`
            Response:
                `{
                    "userId": 1,
                    "status": "active",
                    "walletId (id кошелька)": "hjdgf89345jdsfg98235234df34r34r234r",
                    "createdAt": "2022-01-25T15:38:15+03:00",
                    "managerId": 1
                }`
        2. create an order to cash out
            Request:
                `POST /api/account/{id}/{walletId}/cash
                {
                    "amount": 10.00,
                    ...
                }`
            Response:
                `{
                    "status": "Pending",
                    "createdAt": "2022-01-25T15:38:15+03:00"
                    "cashOutDate": "2022-05-25T12:00:00+03:00"
                }`
            Constraints:
                - available only for account and wallet owner
        3. orders:
            Request:
                `GET /api/account/{id}/orders`
            Response:
                `{
                    "orders": [
                        {
                            "amount": 10.00,
                            "status": "Pending",
                            ...
                        }
                    ]
                }`
    wallet:
        1. wallet balance
            Request:
                `GET /api/wallet/{userId}/{walletId}
                {
                    "id": 1,
                    "ballance": 100.00,
                    "status": "active",
                    ...
                }`
            Constraints:
                - available for owner, manager super admin
            Questions:
                - multiple wallets?
    transaction:
        1. send money
            Request:
                `POST /api/transaction/{userId}
                {
                    "from": "walletSender"
                    "to": "walletReceiver",
                    "amount": 10.00
                }`
            Constraints:
                - transfer available only from own wallet
            Questions:
                - transfer status? (If yes, then we need additional API method to get transfer status)
        2. transaction history
            Request:
                `GET /api/transaction/all/?page=1&limit=30
                {
                    "from": "123",
                    "to": "456",
                    "amount": "10.00",
                    "createdAt": "2022-01-25T15:38:15+03:00",
                    ...
                }`
        3. transaction history certain user
            Request:
               `GET /api/transaction/{userId}/all?page=1&limit=30`
    manager:
        1. create and modify user
            Request:
                `POST|PATCH /api/manager/user
                {
                    "name": "Max",
                    "email": "maxim@mail.com",
                    ...
                }`
            Constraints:
                - modifying available for managers that handle their user and super admin
        2. create and modify user account
            Request:
               `POST|PATCH /api/manager/account
                {
                    "userId": "1",
                    "ballance": "10"
                    ...
                }`
            Constraints:
                - modifying available for manager that handle this user and super admin
        3. orders to cash out:
            Request:
                `GET /api/manager/account/orders?page=1&limit=30&orderStatus=Pending|Done|Rejected`
            Response:
                `{
                    {
                        "accountId": 1,
                        "orders": [
                            {
                                "id": 1,
                                "amount": 10.00,
                                "status": "Pending",
                                ...
                            }
                        ],
                    }                    
                }`
            Constraints:
                - managers can see orders only those users, that they handle, but super admin all. 
        4. handle an order to cash out
            Request:
                `POST /api/manager/order/{orderId}
                {
                    "status": "approve|reject|..."
                    ...
                }`
            Notes:
                - handle by setting status with possible values: approve, reject
            Constraints:
                - managers can handle orders only those users, that they handle, but super admin all.
        5. Debet amount from user wallet
            Request:
                `POST /api/manager/wallet/{walletId}/debet
                {
                    "amount": 10.00,
                    "reason": "cashOut|fine|etc..."
                }`
            Notes:
                - During this request system create a transaction and transfer money to manager's wallet 
                - Super admin can debet money from manager's wallets.
                - Managers can debet money only from those users that they handle