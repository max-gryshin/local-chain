## API

# user
 user list.
 - Request: `GET /api/user/all?page=1&limit=30`
 - Response:
```
{
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
}
```
 certain user
 - Request: `GET /api/user/{id}`
 - Response:
```
{
    "name": "max",
    "dob": 1993.11.11,
    "status": "active"
    ...
}
```
 # account:
 user account
 - Request: `GET /api/account/{userId}/`
 - Response:
```
{
    "userId": 1,
    "status": "active",
    "walletId": "hjdgf89345jdsfg98235234df34r34r234r",
    "createdAt": "2022-01-25T15:38:15+03:00",
    "managerId": 1
}
```
 create an order to cash out
 - Request: `POST /api/account/{id}/{walletId}/cash`
```
{
    "amount": 10.00,
    ...
}
```
 - Response:
```
{
    "status": "Pending",
    "createdAt": "2022-01-25T15:38:15+03:00"
    "cashOutDate": "2022-05-25T12:00:00+03:00"
}
```
 - Constraints:
   - available only for account and wallet owner

 # order:
 - Request: `GET /api/order`
 - Response:
```
{
    "orders": [
        {
            "amount": 10.00,
            "status": "Pending",
            ...
        }
    ]
}
```
 - Constraints:
    - available only for order owner and it's manager

 - Request: `GET /api/order/{id}`
 - Response:
```
{
    "orders": [
        {
            "amount": 10.00,
            "status": "Pending",
            ...
        }
    ]
}
```
- Constraints:
    - available only for order owner and it's manager

 # wallet:
 wallet balance
 - Request: `GET /api/wallet/{walletId}`
```
{
    "id": 1,
    "ballance": 100.00,
    "status": "active",
    ...
}
```
 - Constraints:
   - available for owner, manager super admin
 - Questions:
   - multiple wallets?

# transaction:
 send money
 - Request: `POST /api/transaction`
```
{    
    "to": "walletReceiver",
    "amount": 10.00
}
```
 - Constraints:
   - transfer available only from own wallet
 - Questions:
   - transfer status? (If yes, then we need additional API method to get transfer status)
 
 transaction history
 - Request: `GET /api/transaction/all/?page=1&limit=30`
```
{
    "from": "123",
    "to": "456",
    "amount": "10.00",
    "createdAt": "2022-01-25T15:38:15+03:00",
    ...
}
```
 transaction history certain user
 - Request: `GET /api/transaction/{userId}/all?page=1&limit=30`

# manager
 create and modify user
 - Request: `POST|PATCH /api/manager/user`
```
{
    "name": "Max",
    "email": "maxim@mail.com",
    ...
}
```
 - Constraints:
   - modifying available for managers that handle their user and super admin
 
create and modify user account
 - Request: `POST|PATCH /api/manager/account`
```
{
    "userId": "1",
    "ballance": "10"
    ...
}
```
 - Constraints:
   - modifying available for manager that handle this user and super admin

orders to cash out:
 - Request: `GET /api/manager/order?page=1&limit=30&orderStatus=Pending|Done|Rejected`
 - Response:
```
{
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
}
```
 - Constraints:
   - managers can see orders only those users, that they handle, but super admin all. 

handle an order to cash out
 - Request: `POST /api/manager/order/{orderId}`
```
{
    "status": "approve|reject|..."
    ...
}
```
- Notes:
   - handle by setting status with possible values: approve, reject
- Constraints:
  - managers can handle orders only those users, that they handle, but super admin all.

Debet amount from user wallet
- Request: `POST /api/manager/wallet/{walletId}/debet`
```
{
    "amount": 10.00,
    "reason": "cashOut|fine|etc..."
}
```
- Notes:
    - During this request system create a transaction and transfer money to manager's wallet 
    - Super admin can debet money from manager's wallets.
    - Managers can debet money only from those users that they handle