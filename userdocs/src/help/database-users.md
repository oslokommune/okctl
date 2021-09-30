## Setting up a database user in your cluster


###Create a secret to connect to your database
In AWS console, open Parameter Store and create a secret like this:

Name: /okctl/mycluster/myapp/db_password
Type: SecureString
Value: <some-strong-password, consider using uuidgen>
See uuidgen command at https://okctl.io/getting-started/application-addons/


Use okctl forward and execute the following SQL
See https://okctl.io/getting-started/application-addons/#forwarding-traffic-to-the-database-from-a-local-machine

Create role, and set a password. Use the generated password instead of 'changeme':
```sql
create user my_app with password 'changeme';
```


This depends on app usage:
```sql
alter role my_app with inherit; -- note, inherit means inherit privileges, not role stuff (createrole, createdb, etc)
alter role my_app with createrole;
alter role my_app with nocreatedb;
```