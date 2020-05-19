package main

import (
	"encoding/json"
	"fmt"
	"github.com/emirpasic/gods/maps/treemap"
	"github.com/emirpasic/gods/sets/treeset"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type simap map[string]int
type ssimap map[string]simap

const intern = "实习生"
const outsourcing = "外包"
const socialRecruitment = "社招"
const oncampusRecruitment = "校招"
const fullTime = "全职"
const internalSocialRecruitment = "内推社招"

/*
数据结构

{
	"2020-04" : {
					"外包": 5,
					"实习": 3,
					"校招": 4,
					"社招": 6
				}
}

 */

func main()  {
	log.SetFlags(log.Ldate | log.Lshortfile)
	filename := ""
	if len(os.Args) == 2 {
		if os.Args[1] == "-h" || os.Args[1] == "--help" {
			usage()
			return
		}
		filename = os.Args[1]
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print("Read file %v error : %v\n", filename, err.Error())
		return
	}



	result := make(ssimap)
	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		if i == 0 {
			continue
		}

		log.Print(line, "\n")
		line = strings.TrimSpace(line)
		words := strings.Split(line, ",")
		for i, w := range words {
			words[i] = strings.TrimSpace(w)
		}

		// 过滤掉 撤销、审批不通过 的offer
		if words[10] != "已同意" {
			continue
		}

		// 检查第8列，并获取时间字段
		const dateIndex = 8
		log.Printf("words len = %v\n", len(words))
		date :=  words[dateIndex]
		if len(date) != 10 {
			log.Printf("Date field ERROR, it is not %v index. [%v]", dateIndex, date)
			return
		}
		yearMonth := date[0:7]
		log.Printf("Year month = %v\n", yearMonth)


		//检查第四列，并获取offer类型
		offerType := getOfferType(words[4], words[6])

		kv, ok := result[yearMonth]
		if !ok {
			kv = make(simap)
			result[yearMonth] = kv

			kv[intern] = 0
			kv[fullTime] = 0
			kv[oncampusRecruitment] = 0
			kv[outsourcing] = 0
			kv[socialRecruitment] = 0
		}

		kv[offerType] = kv[offerType] + 1


		//检查第三列，渠道来源
		source := words[3]
		if _, ok := kv[source]; !ok {
			kv[source] = 0
		}
		kv[source] = kv[source] + 1

		//技术内推社招offer数量
		if source == "内部渠道" && offerType == socialRecruitment {
			if _, ok := kv[internalSocialRecruitment]; !ok {
				kv[internalSocialRecruitment] = 0
			}
			kv[internalSocialRecruitment] = kv[internalSocialRecruitment] + 1
		}
	}

	for k, v := range result {
		n := v[oncampusRecruitment] + v[socialRecruitment]
		result[k][fullTime] = n
	}

	fmt.Printf("\n\n\n")
	fmt.Print(result)
	fmt.Printf("\n\n\n")
	j, _ := json.MarshalIndent(result, "", "   ")
	fmt.Printf(string(j))
	fmt.Printf("\n\n\n")


	//每个月组成的时间集合
	//2017-09, 2017-10, 2017-11, 2017-12, 2018-01, 2018-02, 2018-03, 2018-04, 2018-05, 2018-06, 2018-07, 2018-08, 2018-09, 2018-10, 2018-11, 2018-12, 2019-01, 2019-02, 2019-03, 2019-04, 2019-05, 2019-06, 2019-07, 2019-08, 2019-09, 2019-10, 2019-11, 2019-12, 2020-01, 2020-02, 2020-03, 2020-04
	dateSet := treeset.NewWithStringComparator()
	for k, _ := range result {
		dateSet.Add(k)
	}
	log.Printf("%v\n\nsprintf=%v\n", dateSet, fmt.Sprintf("%v", dateSet))


	//来源组成的集合，每一个key都会有对应的一条曲线
	//RPO, 人才库, 全职, 内部渠道, 外包, 实习生, 校招, 猎头渠道, 社交渠道, 社招, 线下渠道, 网络渠道
	lineMap := treemap.NewWithStringComparator()
	for _, v := range result {
		for k, _ := range v {
			lineMap.Put(k, "")
		}
	}
	log.Printf("%v\n\nlineMap sprintf=%v\n", lineMap, fmt.Sprintf("%v", lineMap))


	// 组合每一行的Data曲线数据
	for _, k := range dateSet.Values() {
		vs, _ := result[k.(string)]
		//for k, v := range vs {
		//	line, _ := lineMap.Get(k)
		//	nl := line.(string) + "," + fmt.Sprintf("%v", v)
		//	lineMap.Put(k, nl)
		//}
		keys := lineMap.Keys()
		log.Printf("keys=%v\n", keys)
		for _, k := range keys {
			line, ok := lineMap.Get(k.(string))
			if !ok {
				log.Printf("ERROR line=%v k=%v\n", line, k)
			}

			nn := "0"
			v, ok := vs[k.(string)]
			if ok {
				nn = fmt.Sprintf("%v", v)
			}
			nl := line.(string) + ", " + nn
			lineMap.Put(k, nl)
		}
	}

	log.Printf("%v\n\nlineMap sprintf=%v\n", lineMap, fmt.Sprintf("%v", lineMap))


	/*
ChartType = line
Title = Source: WorldClimate.com
SubTitle = Monthly Average Temperature
ValueSuffix = °C
XAxisNumbers = 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12
YAxisText = Temperature (°C)

# The data and the name of the lines
Data|Tokyo = 7.0, 6.9, 9.5, 14.5, 18.2, 21.5, 25.2, 26.5, 23.3, 18.3, 13.9, 9.6
Data|New York = -0.2, 0.8, 5.7, 11.3, 17.0, 22.0, 24.8, 24.1, 20.1, 14.1, 8.6, 2.5
Data|Berlin = -0.9, 0.6, 3.5, 8.4, 13.5, 17.0, 18.6, 17.9, 14.3, 9.0, 3.9, 1.0
Data|London = 3.9, 4.2, 5.7, 8.5, 11.9, 15.2, 17.0, 16.6, 14.2, 10.3, 6.6, 4.8
	*/

	//d := make(map[string]string)

	chartContent := "ChartType = line" + "\n"
	chartContent += "Title = 每月Offer数据\n"
	chartContent += "SubTitle = \n"
	chartContent += "ValueSuffix = \n"
	chartContent += "XAxisNumbers = " + joinToString(dateSet)
	chartContent = strings.Trim(chartContent, ", ")
	chartContent += "\n"
	chartContent += "YAxisText = \n"

	it := lineMap.Iterator()
	for it.Next() {
		key, value := it.Key(), it.Value()
		v := strings.Trim(value.(string), ", \n\t")
		chartContent += "Data|" + key.(string) + "=" + v + "\n"
	}

	log.Printf("\n\n%v\n", chartContent)

	ioutil.WriteFile("offer.chart", []byte(chartContent), 0644)
}

// 参数：
// 岗位名称，例如 后台开发实习生
// 职级，例如 2-1
// return offer类型
func getOfferType(jobName string, rank string) string {
	log.Printf("job=[%v] rank=[%v]\n", jobName, rank)
	if strings.Contains(jobName, intern) {
		return intern
	}

	if strings.Contains(jobName, outsourcing) {
		return outsourcing
	}

	if rank == "0" {
		return intern
	}

	if rank == "1-2" {
		return oncampusRecruitment
	}

	return socialRecruitment
}

func joinToString(dateSet* treeset.Set) string {
	r := ""
	vv := dateSet.Values()
	for i, v := range vv {
		if i != 0 {
			r += ", "
		}
		s := v.(string)[0:4] + v.(string)[5:]
		r += s
		log.Printf("s=[%v] v=[%v] len(v)=%v\n", s, v, len(v.(string)))
	}

	return r
}


func usage() {
	fmt.Printf("usage : %v offer-filename.csv\n", os.Args[0])
	fmt.Printf("下载offer文件，另存为csv文件。然后运行这个程序，会生成一个 offer.chart 文件。最后 gochart 工具绘图，用浏览器打开即可查看。 http://localhost:8000/ 。\n")
}