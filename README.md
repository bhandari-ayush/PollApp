# PollApp



Postgress img download from registry using calm credential and tag as postgres:pollapp
# Login instead container
psql -U root -d pollapp

# Get list of tables
\dt

# View table schema
\d public.table_name

# View data in table
select * from table public.table_name