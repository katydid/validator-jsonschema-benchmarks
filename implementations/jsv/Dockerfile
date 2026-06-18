FROM elixir:1.20.1-otp-29

RUN apt-get update -qq \
    && apt-get install -qq --no-install-recommends \
      time \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

ENV MIX_ENV=prod

WORKDIR /app
COPY ./mix.exs /app/mix.exs
RUN mix deps.get && mix deps.compile

COPY . /app
RUN mix compile

ENTRYPOINT ["/app/memory-wrapper.sh", "mix", "run", "benchmark.exs"]
CMD []
