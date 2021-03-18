FROM golang:1.16.2-buster AS build

WORKDIR /build
COPY . .

RUN go build -o /bin/api

FROM python:3.8 AS runtime

RUN pip install --no-cache-dir matplotlib joblib sklearn numpy

ARG omw

COPY --from=build /bin/api /bin/api
COPY --from=build /etc/nsswitch.conf /etc/nsswitch.conf
COPY ./UVIndexModel.joblib /app/
COPY ./scripts/getUVIndex.py /app/
ARG buildTime_MODELSIZE='default'


LABEL app="sensoplas-api"

ENV PORT="8080" \
    OMW_API_KEY=${omw} \
    MODELSIZE=$buildTime_MODELSIZE


ENTRYPOINT [ "/bin/api", "http" ]