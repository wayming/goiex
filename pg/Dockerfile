FROM centos
ENV container docker


# ARGS
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
# Set to POXIS locale
ENV LC_ALL POSIX 
ENV LANG POSIX 
RUN sed -e "s#LANG=.*#LANG=\"POSIX\"#" /etc/locale.conf > /etc/locale.conf

# Disable the built-in PostgreSQL module:
RUN dnf -qy module disable postgresql

# Install the repository RPM:
RUN dnf install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm

# Install PostgreSQL and tools:
RUN dnf list -y postgresql*
RUN dnf install -y postgresql13-server
RUN dnf install -y net-tools
RUN dnf install -y python3
ENV PATH="/usr/pgsql-13/bin:${PATH}"
RUN alternatives --set python /usr/bin/python3

# Install systemctl alternative at last, in case systemd is overriden by other installations.
ENV SYSTEMCTL_VERSION=1.4.4181
ADD https://github.com/gdraheim/docker-systemctl-replacement/archive/v${SYSTEMCTL_VERSION}.tar.gz .
RUN tar xvf v${SYSTEMCTL_VERSION}.tar.gz docker-systemctl-replacement-${SYSTEMCTL_VERSION}/files/docker/systemctl.py && \
    /bin/rm -f /usr/bin/systemctl && \
    cp docker-systemctl-replacement-${SYSTEMCTL_VERSION}/files/docker/systemctl.py /usr/bin/systemctl

RUN rm -vf /lib/systemd/system/sysinit.target.wants/* \
  ; rm -vf /lib/systemd/system/multi-user.target.wants/* \
  ; rm -vf /etc/systemd/system/*.wants/* \
  ; rm -vf /lib/systemd/system/local-fs.target.wants/* \
  ; rm -vf /lib/systemd/system/sockets.target.wants/*udev* \
  ; rm -vf /lib/systemd/system/sockets.target.wants/*initctl* \
  ; rm -vf /lib/systemd/system/basic.target.wants/*

# DB setup
# Updates postgres configuration parameters so that it can be started in the container
RUN /usr/pgsql-13/bin/postgresql-13-setup initdb
RUN sed -e "/#listen_addresses/a listen_addresses = '*'" /var/lib/pgsql/13/data/postgresql.conf > /var/lib/pgsql/13/data/postgresql.conf.tmp \
  ; mv /var/lib/pgsql/13/data/postgresql.conf.tmp /var/lib/pgsql/13/data/postgresql.conf
RUN grep 'listen_addresses\|unix_socket_directories' /var/lib/pgsql/13/data/postgresql.conf
RUN echo 'host all all 0.0.0.0/0 md5' >> /var/lib/pgsql/13/data/pg_hba.conf
# # DB create
COPY dbsetup.sh /opt
RUN  chmod +x /opt/dbsetup.sh

USER postgres
RUN  /opt/dbsetup.sh ${DB_USER} ${DB_PASSWORD} ${DB_NAME}

USER root
RUN  whoami
RUN systemctl enable postgresql-13
CMD /usr/bin/systemctl
