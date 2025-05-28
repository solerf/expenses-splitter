FROM alpine/java:21-jdk AS build

RUN apk update; apk add make

RUN java -version

RUN mkdir app-build

WORKDIR /app-build

COPY ./src ./src
COPY ./.mvn ./.mvn
COPY ./mvnw ./config.yml ./pom.xml ./Makefile ./

RUN make build

FROM alpine/java:21-jre AS final

RUN apk update; apk add bash

RUN java -version

RUN mkdir -p /bill-splitter
COPY --from=build /app-build/target/billsplitter-1.0-SNAPSHOT.jar /bill-splitter
COPY --from=build /app-build/config.yml /bill-splitter

WORKDIR /bill-splitter
EXPOSE 8000
CMD [ "java", "-jar", "billsplitter-1.0-SNAPSHOT.jar", "server", "config.yml" ]
