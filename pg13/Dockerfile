FROM centos
ENV container docker


# ARGS
ARG DB_USER
ARG DB_PASSWORD
ARG DB_NAME
ARG DB_VERSION=13
ARG PG_WORKDIR=/opt/pg
# Set to POXIS locale
ENV LC_ALL POSIX 
ENV LANG POSIX 
RUN sed -e "s#LANG=.*#LANG=\"POSIX\"#" /etc/locale.conf > /etc/locale.conf

# Disable the built-in PostgreSQL module:
RUN dnf -qy module disable postgresql

# Install the repository RPM:
RUN dnf install -y https://download.postgresql.org/pub/repos/yum/reporpms/EL-8-x86_64/pgdg-redhat-repo-latest.noarch.rpm

# Install PostgreSQL and tools:
RUN dnf install -y postgresql${DB_VERSION}-server
RUN dnf install -y gcc python3 platform-python-devel.x86_64 libpq-devel.x86_64
RUN dnf install -y net-tools gettext
RUN pip3 install Pyrseas
ENV PATH="/usr/pgsql-${DB_VERSION}/bin:${PATH}"
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
RUN /usr/pgsql-${DB_VERSION}/bin/postgresql-${DB_VERSION}-setup initdb
RUN sed -e "/#listen_addresses/a listen_addresses = '*'" /var/lib/pgsql/${DB_VERSION}/data/postgresql.conf > /var/lib/pgsql/${DB_VERSION}/data/postgresql.conf.tmp \
  ; mv /var/lib/pgsql/${DB_VERSION}/data/postgresql.conf.tmp /var/lib/pgsql/${DB_VERSION}/data/postgresql.conf
RUN grep 'listen_addresses\|unix_socket_directories' /var/lib/pgsql/${DB_VERSION}/data/postgresql.conf
RUN cp /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf  /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf.backup
RUN sed -e "s#\(host *all *all *127.0.0.1/32 *\)ident#\1md5#" /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf > /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf.tmp \
  ; mv /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf.tmp /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf
RUN echo 'host all all 0.0.0.0/0 md5' >> /var/lib/pgsql/${DB_VERSION}/data/pg_hba.conf
# # DB create
RUN  mkdir ${PG_WORKDIR}
RUN  chown postgres ${PG_WORKDIR}
RUN  chgrp postgres ${PG_WORKDIR}
COPY dbsetup.sh ${PG_WORKDIR}
COPY schema/*.yaml.template ${PG_WORKDIR}
RUN  chmod +x ${PG_WORKDIR}/dbsetup.sh

USER postgres
WORKDIR ${PG_WORKDIR}
RUN  sh dbsetup.sh ${DB_USER} ${DB_PASSWORD} ${DB_NAME} ${DB_VERSION}

USER root
RUN  whoami
RUN systemctl enable postgresql-${DB_VERSION}
CMD /usr/bin/systemctl
