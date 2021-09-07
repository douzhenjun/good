'''
@Description:
@version:
@Company: iwhalecloud
@Author: Zhu.Jie
@Date: 2020-04-26 14:55:27
LastEditors: Zhu.Jie
LastEditTime: 2021-03-11 15:17:20
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
        self.menu_t, self.comp_t, self.auth_t, = "menu", "component", "authority"
        self.menu_comp_t, self.comp_auth_t = "menu_component", "component_authority"
        self.menu_c = ["name", "icon", "desc", 'parent_id', 'status', 'desc_en','label']
        self.comp_c = ["name", "priority_id", "desc", "status", 'desc_en']
        self.auth_c = ["name", "priority_id", "api_path", "method", "desc", "selectable", 'desc_en']
        self.menu_comp_c = ["menu_id", "component_id"]
        self.comp_auth_c = ["component_id", "authority_id"]
        self.priority = dict(list=1, detail=2, curd=3, enable_disable=4)

    def _insert_data(self, table, columns, values, return_flag=False):
        # if return_flag:
        #     sql = "select id from " + table + " where " + columns[0] + "=%s"
        #     self.cursor.execute(sql, (values[0],))
        #     row = self.cursor.fetchone()
        #     if row:
        #         return row[0]
        # else:
        #     sql = "select id from " + table + " where " + columns[0] + "=%s and " + columns[1] + "=%s"
        #     self.cursor.execute(sql, (values[0],values[1]))
        #     row = self.cursor.fetchone()
        #     if row:
        #         return
        sql, placeholder = "insert into " + table + "(", ""
        for item in columns:
            item = "\"" + item + "\"" if item == "desc" else item
            sql += item + ","
            placeholder += "%s,"
        sql = sql[:-1] + ") values (" + placeholder[:-1] + ")"
        self.cursor.execute(sql, values)
        self.conn.commit()
        # 返回记录主键
        if return_flag:
            sql = "select id from " + table + " where " + columns[0] + "=%s order by id desc"
            self.cursor.execute(sql, (values[0],))
            row = self.cursor.fetchone()
            return row[0]
        else:
            return None

    def main(self):
        for root, dirs, files in os.walk(sys.path[0]):
            for item in files:
                if 'json' in item:
                    filename = sys.path[0] + "/" + item
                    f = open(filename, 'r', encoding='utf-8')
                    data = json.load(f)
                    self._init_menu(data)

    def _init_menu(self, data):
        # insert 根节点
        root_menu = self._insert_data(self.menu_t, self.menu_c, [data["name"], data["icon"], data["desc"], None, True, data["descEn"],data["label"]], True)
        for item in data["subMenu"]:
            if "component" in item:
                v_data = [item["name"], item["icon"], item["desc"], root_menu, True, item['descEn'],None]
                # insert 二级页节点菜单
                menu = self._insert_data(self.menu_t, self.menu_c, v_data, True)
                self._init_comp(item["component"], menu)
            else:
                v_data = [item["name"], item["icon"], item["desc"], root_menu, True, item['descEn'],None]
                menu = self._insert_data(self.menu_t, self.menu_c, v_data, True)
                for item_sub_menu in item["subMenu"]:
                    v_data = [item_sub_menu["name"], item_sub_menu["icon"], item_sub_menu["desc"], menu, True, item_sub_menu['descEn'],None]
                    # insert 二级页节点菜单
                    sub_menu = self._insert_data(self.menu_t, self.menu_c, v_data, True)
                    self._init_comp(item_sub_menu["component"], sub_menu)

    def _init_comp(self, comp_data, menu):
        for item_comp in comp_data:
            v_data = [item_comp["name"], self.priority[item_comp["priority"]], item_comp["desc"], True, item_comp['descEn']]
            comp = self._insert_data(self.comp_t, self.comp_c, v_data, True)
            # 菜单和组件关系
            self._insert_data(self.menu_comp_t, self.menu_comp_c, [menu, comp])
            # 权限
            if item_comp["permission"]:
                perm = item_comp["permission"][0]
                v_data = [perm["name"], self.priority[perm["priority"]], perm["path"], perm["method"], perm["desc"], True, perm['descEn']]
                auth = self._insert_data(self.auth_t, self.auth_c, v_data, True)
                # 组件和权限关系
                self._insert_data(self.comp_auth_t, self.comp_auth_c, [comp, auth])

if __name__ == '__main__':
    if sys.argv[1] and sys.argv[2] and sys.argv[3]:
        init = InitData(sys.argv[1], sys.argv[2], sys.argv[3])
        init.main()
    else:
        print('Please input parameter: host, port, password')
