set -e

PG_USER="root"
PG_PASSWORD="postgres"
PG_DATABASE="pollapp"

sudo yum update -y

echo "### Checking and installing epel-release if not already installed"
if ! rpm -q epel-release &> /dev/null; then
    sudo yum install -y epel-release
else
    echo "epel-release is already installed"
fi

echo "### Checking and installing wget if not already installed"
if ! command -v wget &> /dev/null; then
    sudo yum install -y wget
else
    echo "wget is already installed"
fi

echo "### Checking and installing firewalld if not already installed"
if ! systemctl is-active --quiet firewalld; then
    sudo yum install -y firewalld
    sudo systemctl enable --now firewalld
else
    echo "firewalld is already installed and running"
fi

echo "### starting db server installation"

echo "### Checking and installing pgdg-redhat-repo if not already installed"
if ! rpm -q pgdg-redhat-repo &> /dev/null; then
    sudo yum install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-7-x86_64/pgdg-redhat-repo-latest.noarch.rpm
else
    echo "pgdg-redhat-repo is already installed"
fi

echo "### Checking and installing PostgreSQL server and contrib if not already installed"
if ! rpm -q postgresql-server &> /dev/null; then
    sudo yum install -y postgresql-server
else
    echo "postgresql-server is already installed"
fi

if ! rpm -q postgresql-contrib &> /dev/null; then
    sudo yum install -y postgresql-contrib
else
    echo "postgresql-contrib is already installed"
fi

sudo postgresql-setup initdb
sudo systemctl enable postgresql
sudo systemctl start postgresql

sudo -u postgres psql -tc "SELECT 1 FROM pg_roles WHERE rolname = '$PG_USER'" | grep -q 1 || \
sudo -u postgres psql -c "CREATE USER $PG_USER WITH PASSWORD '$PG_PASSWORD';"
sudo -u postgres psql -tc "SELECT 1 FROM pg_database WHERE datname = '$PG_DATABASE'" | grep -q 1 || \
sudo -u postgres psql -c "CREATE DATABASE $PG_DATABASE OWNER $PG_USER;"
sudo -u postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $PG_DATABASE TO $PG_USER;"
sudo -u postgres psql -d $PG_DATABASE -c "ALTER ROLE $PG_USER SET search_path TO public;"

TABLE_SQL="SET ROLE $PG_USER;"

echo "$TABLE_SQL" | sudo -u postgres psql -d "$PG_DATABASE" -v ON_ERROR_STOP=1

PG_CONF="/var/lib/pgsql/data/postgresql.conf"
sudo sed -i "s/^#listen_addresses = .*/listen_addresses = '*'/g" "$PG_CONF"

PG_HBA="/var/lib/pgsql/data/pg_hba.conf"
echo "host    all             all             0.0.0.0/0               md5" | sudo tee -a "$PG_HBA"
sudo sed -i 's/^\(local[[:space:]]\+all[[:space:]]\+all[[:space:]]\+\)peer/\1md5/' /var/lib/pgsql/data/pg_hba.conf
sudo systemctl restart postgresql

echo "### db server installation completed"

sudo firewall-cmd --permanent --add-port=5432/tcp
sudo firewall-cmd --reload
echo "### PostgreSQL installation and configuration complete"
echo "### PostgreSQL is running and listening on port 5432"