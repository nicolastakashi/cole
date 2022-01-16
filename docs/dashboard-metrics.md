# Dashboard Metrics

| Name                        | Type    | Tags                                                                                                                                                                                  |
|-----------------------------|---------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| dashboard_view_total        | Counter | `dashboard_uid`=&lt;dashboard_uid&gt; <br> `org_id`=&lt;org_id&gt; <br> `user_id`=&lt;user_id&gt; <br>`user_name`=&lt;user_name `defult unkown`&gt;                                   |
| dashboard_last_view_seconds | Gauge   | `dashboard_uid`=&lt;dashboard_uid&gt; <br> `org_id`=&lt;org_id&gt;                                                                                                                    |
| dashboard_info              | Gauge   | `dashboard_uid`=&lt;dashboard_uid&gt; <br> `is_stared`=&lt;is_stared&gt; <br> `version`=&lt;version&gt; <br> `schema_version`=&lt;schema_version&gt; <br> `timezone`=&lt;timezone&gt; |