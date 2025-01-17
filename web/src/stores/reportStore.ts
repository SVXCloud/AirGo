import {defineStore} from "pinia";
//report api
import {useReportApi} from "/@/api/report";
import {ElMessage} from "element-plus";

const reportApi = useReportApi()
export const useReportStore = defineStore("reportStore", {
    state: () => ({

        //当前数据库的类型和数据库列表
        dbInfo: {
            db_type: '',
            database_list: [],
        },
        //选中的数据库，库表，用来请求获取数据库的数据表的所有字段名,类型值 请求参数
        checkedDbInfo: {
            database: '',
            table_name: '',
        },
        //高级查询的条件参数
        reportTable: {
            table_name: '',
            field_params_list: [] as FieldTable[],    //搜索条件列表 {field: '', field_chinese_name: '', field_type: '', condition: '=', condition_value: '',}
            pagination_params: {},//分页参数

        },

        //字段信息
        fieldData:{
            field_list:[],
            field_chinese_name_list:{} as {[key:string]: any;},
            field_type_list:{} as {[key:string]: any;},
        },

    }),
    actions: {
        // 获取数据库的所有数据库名
        async getDB(params?: object) {
            const res = await reportApi.getDBApi()
            if (res.code === 0) {
                this.dbInfo = res.data
            }
        },
        // 获取数据库的所有表名，参数：{"database":"xxxx}
        async getTables(params?: object) {
            const res = await reportApi.getTablesApi(params)

        },
        // 获取字段名,类型值
        // 参数：{"table_name":"xxx}
        async getColumn(params?: object) {
            const res = await reportApi.getColumnApi(params)
            this.fieldData=res.data


        },

    }

})