FROM xanderflood/plaid-ui:local

COPY ./start.sh ./start.sh

CMD ["./start.sh"]
