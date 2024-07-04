# API specification
create user:
POST /user
{
    id: string, primary key
    secret: string, hash of password
    salt: string
    
}

Protected endpoints: use authorization header
GET /user/:id
PATCH /user/:id
    do not allow changing id

