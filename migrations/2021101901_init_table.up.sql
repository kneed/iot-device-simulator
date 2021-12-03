
CREATE TABLE IF NOT EXISTS "public"."device" (
                                   "id" int8 NOT NULL,
                                   "created_at" timestamptz(6),
                                   "updated_at" timestamptz(6),
                                   "deleted_at" timestamptz(6),
                                   "name" varchar COLLATE "pg_catalog"."default",
                                   "type" varchar COLLATE "pg_catalog"."default",
                                   "server_ip" varchar COLLATE "pg_catalog"."default",
                                   "server_port" varchar COLLATE "pg_catalog"."default",
                                   "state" varchar COLLATE "pg_catalog"."default"
)
;
COMMENT ON COLUMN "public"."device"."name" IS '设备名字';
COMMENT ON COLUMN "public"."device"."type" IS '设备类型';
COMMENT ON COLUMN "public"."device"."server_ip" IS 'MQTT服务ip';
COMMENT ON COLUMN "public"."device"."server_port" IS 'MQTT服务端口';
COMMENT ON COLUMN "public"."device"."state" IS '设备状态';

-- ----------------------------
-- Records of device
-- ----------------------------
INSERT INTO "public"."device" VALUES (2, '2021-11-30 15:11:07.295682+08', '2021-12-03 09:02:15.350256+08', NULL, '测试模拟设备2', 'T1', '172.19.176.173', '1883', 'error');
INSERT INTO "public"."device" VALUES (1, '2021-11-25 17:45:51.757963+08', '2021-12-03 09:02:17.012829+08', NULL, 'keal', 'T1', '172.19.176.173', '1883', 'error');

-- ----------------------------
-- Table structure for protocol
-- ----------------------------

CREATE TABLE IF NOT EXISTS "public"."protocol" (
                                     "id" int8 NOT NULL,
                                     "created_at" timestamptz(6),
                                     "updated_at" timestamptz(6),
                                     "deleted_at" timestamptz(6),
                                     "device_id" int8,
                                     "name" varchar COLLATE "pg_catalog"."default",
                                     "content" jsonb,
                                     "qos" int2,
                                     "type" int2,
                                     "sub_topic" varchar COLLATE "pg_catalog"."default",
                                     "pub_topic" varchar COLLATE "pg_catalog"."default",
                                     "strategy" jsonb
)
;
COMMENT ON COLUMN "public"."protocol"."device_id" IS '设备id';
COMMENT ON COLUMN "public"."protocol"."name" IS '协议名字';
COMMENT ON COLUMN "public"."protocol"."content" IS '协议内容';
COMMENT ON COLUMN "public"."protocol"."qos" IS 'MQTT消息qos等级';
COMMENT ON COLUMN "public"."protocol"."type" IS '响应类型: 1-自发 2-响应';
COMMENT ON COLUMN "public"."protocol"."sub_topic" IS '订阅MQTT topic';
COMMENT ON COLUMN "public"."protocol"."pub_topic" IS '发布MQTT topic';
COMMENT ON COLUMN "public"."protocol"."strategy" IS '发送策略';

-- ----------------------------
-- Records of protocol
-- ----------------------------
INSERT INTO "public"."protocol" VALUES (2, '2021-11-30 14:03:02.766246+08', '2021-11-30 14:03:02.766246+08', NULL, 1, '订阅测试', '{"purpose": "test_subscribe"}', 1, 1, 'v1/simulator/subscribe/test/1', 'v1/simulator/publish/test/1', '{"duration": 0, "interval": 20}');
INSERT INTO "public"."protocol" VALUES (4, '2021-12-01 10:38:25.37411+08', '2021-12-01 10:38:25.37411+08', NULL, 2, '订阅测试', '{"purpose": "test_subscribe"}', 1, 1, 'v1/simulator/subscribe/test/2', 'v1/simulator/publish/test/2', '{"duration": 0, "interval": 20}');
INSERT INTO "public"."protocol" VALUES (3, '2021-12-01 10:38:10.984667+08', '2021-12-01 10:38:10.984667+08', NULL, 2, '发送测试', '{"purpose": "test2"}', 1, 0, '', 'v1/simulator/publish/test/2', '{"duration": 0, "interval": 20}');
INSERT INTO "public"."protocol" VALUES (1, '2021-11-30 10:55:30.473733+08', '2021-11-30 10:55:30.473733+08', NULL, 1, '发送测试', '{"purpose": "test1"}', 0, 0, '', 'v1/simulator/publish/test/1', '{"duration": 0, "interval": 10}');

-- ----------------------------
-- Indexes structure for table device
-- ----------------------------
CREATE INDEX "idx_deivce_deleted_at" ON "public"."device" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table device
-- ----------------------------
ALTER TABLE "public"."device" ADD CONSTRAINT "deivce_pkey" PRIMARY KEY ("id");

-- ----------------------------
-- Indexes structure for table protocol
-- ----------------------------
CREATE INDEX "idx_protocol_deleted_at" ON "public"."protocol" USING btree (
  "deleted_at" "pg_catalog"."timestamptz_ops" ASC NULLS LAST
);
CREATE INDEX "idx_protocol_device_id" ON "public"."protocol" USING btree (
  "device_id" "pg_catalog"."int8_ops" ASC NULLS LAST
);

-- ----------------------------
-- Primary Key structure for table protocol
-- ----------------------------
ALTER TABLE "public"."protocol" ADD CONSTRAINT "protocol_pkey" PRIMARY KEY ("id");
