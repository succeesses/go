select currency_id, auditor, sum(count1) as total_num, time from
((SELECT currency_id, auditor, date_format(add_time, '%Y-%m-%d') as time, sum(num) as count1 from iwala_deposit group by currency_id)
union all (SELECT currency_id, auditor, date_format(add_time, '%Y-%m-%d') as time, sum(num) as count1 from iwala_withdraw group by currency_id))
as total group by total.currency_id having total.auditor <> 0;
