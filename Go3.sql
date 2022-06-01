select c.time, c.currency_name, sum(c.total) as Total, c.currency_trade_id from
(select date_format(FROM_UNIXTIME(a.add_time), '%H') as time, b.currency_name, a.currency_trade_id, sum(a.money) as total
from iwala_trade as a inner join iwala_currency as b on a.currency_id = b.currency_id
where a.currency_trade_id = 104 group by time order by time)  as c group by c.currency_name;