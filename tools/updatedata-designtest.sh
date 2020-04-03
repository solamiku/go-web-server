#!/bin/bash
# tj921-策划测试服专用 更新策划数据 取db为测试用数值目录下的db
cd /home/stararc/rangers-local/local-server/pre-design/gamesrv/ 
rm -rf data/design.db
svn up
svn export --force svn://tj901/taojin/trunk/v5/design/策划文档/测试用数值/design.db data/design.db
