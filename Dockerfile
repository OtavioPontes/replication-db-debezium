FROM quay.io/debezium/connect:1.9

ENV CONNECTOR_VERSION=10.8.4
ENV CONNECT_PLUGIN_DIR=/kafka/connect

RUN curl -sL https://hub-downloads.confluent.io/api/plugins/confluentinc/kafka-connect-jdbc/versions/${CONNECTOR_VERSION}/confluentinc-kafka-connect-jdbc-${CONNECTOR_VERSION}.zip \
    -o /tmp/confluentinc-kafka-connect-jdbc-${CONNECTOR_VERSION}.zip && \
    unzip -o /tmp/confluentinc-kafka-connect-jdbc-${CONNECTOR_VERSION}.zip -d ${CONNECT_PLUGIN_DIR} && \
    rm /tmp/confluentinc-kafka-connect-jdbc-${CONNECTOR_VERSION}.zip