FROM maven:3.9.7-eclipse-temurin-21@sha256:3b5f7c15b1a16d3fdf09e6883cde602e4a5406cf5bdf6b251b8ac5c510219311 AS builder

WORKDIR /app
 
COPY .mvn/ .mvn
COPY mvnw pom.xml ./
RUN ./mvnw dependency:go-offline
 
COPY src ./src
 
RUN ./mvnw clean install


FROM eclipse-temurin:21.0.5_11-jre-alpine@sha256:4300bfe1e11f3dfc3e3512f39939f9093cf18d0e581d1ab1ccd0512f32fe33f0

COPY --from=builder /app/target/*.jar app.jar

ENTRYPOINT ["java", "-jar", "/app.jar"]
