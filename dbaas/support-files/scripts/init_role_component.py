'''
@Description:
@version:
@Company: iwhalecloud
@Author: Zhu.Jie
@Date: 2020-04-26 14:55:27
@LastEditors: Zhu.Jie
@LastEditTime: 2020-05-09 15:01:41
'''
import os
import psycopg2
import json
import sys


class InitData():
    def __init__(self, host, port, password):
        self.conn = psycopg2.connect(
            host=host, port=port, database='capricorn', user='postgres', password=password)
        self.cursor = self.conn.cursor()
        self.component_name = ["guide.list", "home.list", "dbaas.overview.base",
                               "dbaas.pod.list", "dbaas.storage.init.add",
                               "dbaas.overview.chart", "dbaas.overview.threed",
                               "dbaas.image.init.add",
                               "dbaas.cluster.add", "dbaas.cluster.update",
                               "dbaas.cluster.delete", "dbaas.cluster.operate",
                               "dbaas.pod.detail", "dbaas.pod.switch", "dbaas.pod.log",
                               "dbaas.alarm.log.list", "dbaas.alarm.level.list"]
        self.sql_list = ["""
                SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.cluster.list';
                """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.user.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.cluster')) AND "name"='dbaas.cluster.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.parameter.update';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.host.tag';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.host.operator';

                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.image.type.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.parameter.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.overview')) AND "name"='dbaas.storage.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.cluster')) AND "name"='dbaas.image.param';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.cluster')) AND "name"='dbaas.image.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.cluster')) AND "name"='dbaas.selfUser.list';
                       """,
                         """
                       SELECT id FROM component WHERE id in (SELECT component_id FROM menu_component WHERE menu_id=(SELECT id FROM menu WHERE "name"='dbaas.storage')) AND "name"='dbaas.storage.list';
                       """,
                         ]
        self.comt_id_list = []

    def init_role(self):
        init_role_sql = """
INSERT INTO "public"."role"("name", "desc", "is_preseted") VALUES ('dbaas', 'dbaas', 't');
"""
        self.cursor.execute(init_role_sql)
        self.conn.commit()

    def get_compent_name(self):
        sql = """
        SELECT name FROM component WHERE id in (SELECT component_id FROM role_component WHERE role_id = (SELECT id FROM "role" WHERE "name"='dbaas'));
        """
        self.cursor.execute(sql)
        result = self.cursor.fetchall()
        print(result)

    def get_component(self):
        for name in self.component_name:
            sql = """
        SELECT id FROM component WHERE "name"='%s';
                   """ % name
            self.cursor.execute(sql)
            result_list = self.cursor.fetchall()
            print(result_list)
            if result_list:
                for result in result_list:
                    not_root_comp_id = result[0]
                    self.comt_id_list.append(not_root_comp_id)
                    print(not_root_comp_id)

        for sql in self.sql_list:
            self.cursor.execute(sql)
            result_list = self.cursor.fetchall()
            print(result_list)
            if result_list:
                for result in result_list:
                    not_root_comp_id = result[0]
                    self.comt_id_list.append(not_root_comp_id)
                    print(not_root_comp_id)

    def insert_role_component(self):
        comt_ids = list(set(self.comt_id_list))
        for com_id in comt_ids:
            insert_sql = """
            INSERT INTO "public"."role_component"("role_id", "component_id") VALUES ((SELECT id FROM "role" WHERE "name"='dbaas'), %s);

            """ % com_id
            print(insert_sql)
            self.cursor.execute(insert_sql)
            self.conn.commit()

    def main(self):
        self.init_role()
        self.get_component()
        self.insert_role_component()


if __name__ == '__main__':
    if sys.argv[1] and sys.argv[2] and sys.argv[3]:
        init = InitData(sys.argv[1], sys.argv[2], sys.argv[3])
        init.main()
    else:
        print('Please input parameter: host, port, password')
