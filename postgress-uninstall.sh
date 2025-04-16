set -e

echo "### Starting PostgreSQL uninstallation process"

echo "### Removing PostgreSQL packages"
sudo yum remove -y postgresql-server postgresql-contrib || echo "Failed to remove PostgreSQL packages"

echo "### Cleaning up PostgreSQL data and log directories"
sudo rm -rf /var/lib/pgsql /var/log/postgresql || echo "Failed to remove PostgreSQL directories"

echo "### PostgreSQL uninstallation complete"