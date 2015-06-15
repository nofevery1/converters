package models

import (
	_ "bytes"
	"encoding/json"
	_ "fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
	"gopkg.in/mgo.v2/bson"

	"nofe/db"
)

type Article struct {
	Id           bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	ImportedBy   string        `bson:"importedBy,omitempty" json:"importedBy,omitempty"`
	ImportedDate time.Time     `bson:"importedDate,omitempty" json:"importedDate,omitempty"`
	Type         string        `bson:"type,omitempty" json:"type,omitempty"`
	Journal      Journal       `bson:"journal,omitempty" json:"journal,omitempty"`
	Pmid         string        `bson:"pmid,omitempty" json:"pmid,omitempty"`
	Pmc          string        `bson:"pmc,omitempty" json:"pmc,omitempty"`
	Doi          string        `bson:"doi,omitempty" json:"doi,omitempty"`
	PublisherId  string        `bson:"publisherId,omitempty" json:"publisherId,omitempty"`
	Categories   []Category    `bson:"categories,omitempty" json:"categories,omitempty"`
	Titles       []string      `bson:"titles,omitempty" json:"titles,omitempty"`
	Contributors []Contributor `bson:"contributors,omitempty" json:"contributors,omitempty"`
	Aff          Aff           `bson:"aff,omitempty" json:"aff,omitempty"`
	AuthorNotes  []AuthorNote  `bson:"authorNotes,omitempty" json:"authorNotes,omitempty"`
	Ppub         Date          `bson:"ppub,omitempty" json:"ppub,omitempty"`
	Epub         Date          `bson:"epub,omitempty" json:"epub,omitempty"`
	PmcRelease   Date          `bson:"pmcRelease,omitempty" json:"pmcRelease,omitempty"`
	Volume       string        `bson:"volume,omitempty" json:"volume,omitempty"`
	Issue        string        `bson:"issue,omitempty" json:"issue,omitempty"`
	Fpage        string        `bson:"fPage,omitempty" json:"fPage,omitempty"`
	Lpage        string        `bson:"lPage,omitempty" json:"lPage,omitempty"`
	History      History       `bson:"history,omitempty" json:"history,omitempty"`
	Permissions  Permissions   `bson:"permissions,omitempty" json:"permissions,omitempty"`
	Abstract     []Node        `bson:"abstract,omitempty" json:"abstract,omitempty"`
	PageCount    int           `bson:"pageCount,omitempty" json:"pageCount,omitempty"`
	Metas        []Meta        `bson:"metas,omitempty" json:"metas,omitempty"`
	Body         []Node        `bson:"body,omitempty" json:"body,omitempty"`
	Ack          []Node        `bson:"ack,omitempty" json:"ack,omitempty"`
	Refs         Refs          `bson:"refs,omitempty" json:"refs,omitempty"`
}

type Journal struct {
	NlmTa       string   `bson:"nlmTa,omitempty" json:"nlmTa,omitempty"`
	IsoAbbrev   string   `bson:"isoAbbrev,omitempty" json:"isoAbbrev,omitempty"`
	PublisherId string   `bson:"publisherId,omitempty" json:"publisherId,omitempty"`
	Hwp         string   `bson:"hwp,omitempty" json:"hwp,omitempty"`
	Titles      []string `bson:"titles,omitempty" json:"titles,omitempty"`
	Ppub        string   `bson:"ppub,omitempty" json:"ppub,omitempty"`
	Epub        string   `bson:"epub,omitempty" json:"epub,omitempty"`
}

type Category struct {
	Group   string `bson:"group,omitempty" json:"group,omitempty"`
	Subject string `bson:"subject,omitempty" json:"subject,omitempty"`
}

type Contributor struct {
	Type       string `bson:"type,omitempty" json:"type,omitempty"`
	Surname    string `bson:"surname,omitempty" json:"surname,omitempty"`
	GivenNames string `bson:"givenNames,omitempty" json:"givenNames,omitempty"`
	UserId     string `bson:"userId,omitempty" json:"userId,omitempty"`
}

type Aff struct {
	Id       string `bson:"id,omitempty" json:"id,omitempty"`
	Children []Node `bson:"children,omitempty" json:"children,omitempty"`
}

type AuthorNote struct {
	CorrespId string `bson:"correspId,omitempty" json:"correspId,omitempty"`
	Children  []Node `bson:"children,omitempty" json:"children,omitempty"`
}

type History struct {
	Received Date `bson:"received,omitempty" json:"received,omitempty"`
	RevRecd  Date `bson:"revRecd,omitempty" json:"revRecd,omitempty"`
	Accepted Date `bson:"accepted,omitempty" json:"accepted,omitempty"`
}

type Permissions struct {
	CopyrightStatement string   `bson:"copyrightStatement,omitempty" json:"copyrightStatement,omitempty"`
	CopyrightYear      string   `bson:"copyrightYear,omitempty" json:"copyrightYear,omitempty"`
	Licenses           []string `bson:"licenses,omitempty" json:"licenses,omitempty"`
}

type Date struct {
	Day   string `bson:"day,omitempty" json:"day,omitempty"`
	Month string `bson:"month,omitempty" json:"month,omitempty"`
	Year  string `bson:"year,omitempty" json:"year,omitempty"`
}

type Node struct {
	Type      string            `bson:"type,omitempty" json:"type,omitempty"`         //"text" or "tag"
	Tag       string            `bson:"tag,omitempty" json:"tag,omitempty"`           //only for tag types
	Props     map[string]string `bson:"props,omitempty" json:"props,omitempty"`       //only for tag types
	Body      string            `bson:"body,omitempty" json:"body,omitempty"`         //only for text types
	Children  []Node            `bson:"children,omitempty" json:"children,omitempty"` //only for tag types
	Sentences []Sentence        `bson:"sentences,omitempty" json:"sentences,omitempty"`
}

type Sentence struct {
	Start int `bson:"start" json:"start"`
	End   int `bson:"end" json:"end"`
}

type Meta struct {
	Name  string `bson:"name,omitempty" json:"name,omitempty"`
	Value string `bson:"value,omitempty" json:"value,omitempty"`
}

type Refs struct {
	Title string `bson:"refs,omitempty" json:"refs,omitempty"`
	List  []Ref  `bson:"list,omitempty" json:"list,omitempty"`
}

type Ref struct {
	Id     string  `bson:"ref,omitempty" json:"ref,omitempty"`
	Label  string  `bson:"label,omitempty" json:"label,omitempty"`
	Groups []Group `bson:"groups,omitempty" json:"groups,omitempty"`
	Title  string  `bson:"title,omitempty" json:"title,omitempty"`
	Source string  `bson:"source,omitempty" json:"source,omitempty"`
	Year   string  `bson:"year,omitempty" json:"year,omitempty"`
	Volume string  `bson:"volume,omitempty" json:"volume,omitempty"`
	Fpage  string  `bson:"fPage,omitempty" json:"fPage,omitempty"`
	Lpage  string  `bson:"lPage,omitempty" json:"lPage,omitempty"`
	Pmid   string  `bson:"pmid,omitempty" json:"pmid,omitempty"`
}

type Group struct {
	Type  string              `bson:"type,omitempty" json:"type,omitempty"`
	Names []map[string]string `bson:"names,omitempty" json:"names,omitempty"`
}

type EutilsSearch struct {
	ESearchResult EutilsSearchResult `bson:"esearchresult,omitempty" json:"esearchresult,omitempty"`
}

type EutilsSearchResult struct {
	IdList []string `bson:"idlist,omitempty" json:"idlist,omitempty"`
}

func (this *Article) DecodeJSON(reader io.ReadCloser) error {
	defer reader.Close()
	decoder := json.NewDecoder(reader)
	return decoder.Decode(&this)
}

func ArticleParseNodes(nodes []Node, depth int, sentences []Sentence, pmc string) (output string) {

	/*
		Type      string            `bson:"type,omitempty" json:"type,omitempty"`         //"text" or "tag"
		Tag       string            `bson:"tag,omitempty" json:"tag,omitempty"`           //only for tag types
		Props     map[string]string `bson:"props,omitempty" json:"props,omitempty"`       //only for tag types
		Body      string            `bson:"body,omitempty" json:"body,omitempty"`         //only for text types
		Children  []Node            `bson:"children,omitempty" json:"children,omitempty"` //only for tag types
		Sentences []Sentence
	*/

	depth++
	sentenceInd := 0
	for _, node := range nodes {
		//open node
		if node.Type == "text" {
			//parse sentences
			body := node.Body

			//log.Println("ind:", sentenceInd)

			//check if sentence is started or ended
			bodyEnd := len(body) + sentenceInd
			for i := len(sentences) - 1; i >= 0; i-- {
				sent := sentences[i]

				if sent.End <= bodyEnd && sent.End >= sentenceInd {
					end := sent.End - sentenceInd

					if end != len(body) {
						body = body[0:end] + "</div>" + body[end:len(body)]
					} else {
						body = body[0:end] + "</div>"
					}

				}
				if sent.Start <= bodyEnd && sent.Start >= sentenceInd {
					start := sent.Start - sentenceInd
					//log.Println("sent start:", sent, start, len(body), body)
					//log.Println("sent start:", start)

					if start != 0 {
						body = body[0:start] + "<div class=\"ae-sentence\">" + body[start:len(body)]
					} else {
						body = "<div class=\"ae-sentence\">" + body[start:len(body)]
					}

				}
			}

			sentenceInd += len(node.Body)

			//loop backwards through
			//body = body[]+"<div class=\"ae-sentence\">"

			if len(body) > 0 {
				if body[0:1] == ")" {
					body = body[1:len(body)]
				}
				if body[len(body)-1:len(body)] == "(" {
					body = body[:len(body)-1]
				}
			}

			output += body
		}
		if node.Type == "tag" {
			tag := node.Tag
			props := map[string]string{}
			if tag == "title" {
				if depth < 3 {
					tag = "h2"
				} else if depth < 4 {
					//log.Println("depth:", depth)
					tag = "h3"
				} else {
					tag = "h4"
				}
			}
			if tag == "sec" {
				tag = "div"
				props["class"] = "ae-paragraph"
			}
			if tag == "p" {
				tag = "div"
				props["class"] = "ae-paragraph"
			}
			if tag == "xref" {
				tag = "a"
				props["href"] = "javascript:"
				props["data-addition-id"] = "citation"
				props["class"] = "article-additional"
			}
			if tag == "ext-link" {
				tag = "a"
				props["href"] = "javascript:"
				props["data-addition-id"] = "citation"
				props["class"] = "long-word"
			}
			if tag == "inline-formula" {
				//log.Println("formula:", node.Props)
				tag = "span"
			}
			if tag == "disp-formula" {
				//log.Println("formula:", node.Props)
				tag = "div"
				props["class"] = "mobile-small ae-paragraph"
			}
			if tag == "inline-graphic" {
				//log.Println("graphic:", node.Props)
				tag = "img"
				props["src"] = "http://www.ncbi.nlm.nih.gov/pmc/articles/PMC" + pmc + "/bin/" + node.Props["xlink:href"]
				//http://www.ncbi.nlm.nih.gov/pmc/articles/PMC3592458/bin/gks981i3.jpg
			}
			if tag == "graphic" {
				//log.Println("graphic:", node.Props)
				tag = "img"
				props["src"] = "http://www.ncbi.nlm.nih.gov/pmc/articles/PMC" + pmc + "/bin/" + node.Props["xlink:href"] + ".jpg"
				//http://www.ncbi.nlm.nih.gov/pmc/articles/PMC3592458/bin/gks981i3.jpg
			}
			if tag == "fig" {
				//is a figure - need to import figure into side panel
				tag = "fig"
				props["style"] = "display:none;"
			}

			output += "<" + tag
			//write props
			for prop, val := range props {
				output += " " + prop + "=\"" + val + "\""
			}
			output += ">"
			if node.Tag == "xref" {
				output += "["
			}
			output += ArticleParseNodes(node.Children, depth, node.Sentences, pmc)
			if node.Tag == "xref" {
				output += "]"
			}
			//close node
			output += "</" + tag + ">"
		}
	}
	return
}

func ArticleGetByDoi(doi string) (article Article, err error) {
	query := bson.M{
		"doi": doi,
	}

	log.Println("get by doi:", doi)

	err = db.GetCol("articles").Find(query).One(&article)
	if err != nil {
		return
	}
	return
}

func ArticleGetByPmc(pmc string) (article Article, err error) {
	query := bson.M{
		"pmc": pmc,
	}

	log.Println("get by pmc:", pmc)

	err = db.GetCol("articles").Find(query).One(&article)
	if err != nil {
		return
	}
	return
}

func ArticleGetById(id string) (article Article, err error) {
	if !bson.IsObjectIdHex(id) {
		return
	}

	objectId := bson.ObjectIdHex(id)
	query := bson.M{
		"_id": objectId,
	}

	err = db.GetCol("articles").Find(query).One(&article)
	if err != nil {
		return
	}
	return
}

func ArticleImportByDoi(doi string) (article Article, err error) {
	//get doi
	pmc, err := ArticlePmcByDoi(doi)
	if err != nil || pmc == "" {
		return
	}

	//check if pmc already exists in system
	article, err = ArticleGetByPmc(pmc)
	if err != nil && err.Error() != "not found" {
		return
	}

	if article.Pmc == "" {
		article, err = ArticleImportByPmc(pmc)
	}

	return
}

func ArticlePmcByDoi(doi string) (pmc string, err error) {
	resp, err := http.Get("http://eutils.ncbi.nlm.nih.gov/entrez/eutils/esearch.fcgi?retmode=json&db=pmc&term=" + doi)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	var search EutilsSearch
	err = json.Unmarshal(body, &search)
	if err != nil {
		return
	}

	if len(search.ESearchResult.IdList) != 1 {
		return
	}

	if search.ESearchResult.IdList[0] == "" {
		return
	}

	pmc = search.ESearchResult.IdList[0]

	return
}

func ArticleImportByPmc(pmc string) (article Article, err error) {
	resp, err := http.Get("http://eutils.ncbi.nlm.nih.gov/entrez/eutils/efetch.fcgi?db=pmc&name=pubchase&retmode=xml&id=" + pmc)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	//replace tags

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return
	}

	err = article.Parse(doc)
	if err != nil {
		return
	}

	log.Println("article imported")

	if article.Pmc == "" {
		return
	}

	err = db.GetCol("articles").Insert(article)

	return
}

func (this *Article) Parse(n *html.Node) (err error) {
	if n.Data == "article" {
		return this.ParseArticle(n)
		//log.Println(n.Data, n.Attr, level)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.Parse(c)
	}

	return
}

func (this *Article) ParseArticle(n *html.Node) (err error) {

	if n.Data == "journal-meta" {
		return this.ParseJournal(n)
	}
	if n.Data == "article-meta" {
		return this.ParseMeta(n)
	}
	//log.Println("data:", n.Data, n.Namespace, n.DataAtom, n.Type, n.Attr)
	if n.Data == "sec" {

		children, err := this.ParseChildren(n)
		sec := Node{
			Type:     "tag",
			Tag:      n.Data,
			Children: children,
		}
		sec.Props = map[string]string{}
		for _, a := range n.Attr {
			sec.Props[a.Key] = a.Val
		}

		this.Body = append(this.Body, sec)
		return err
	}
	if n.Data == "back" {
		return this.ParseBack(n)
	}

	//log.Println("data:", n.Data, n.Namespace, n.DataAtom, n.Type, n.Attr)

	for _, a := range n.Attr {
		if a.Key == "article-type" {
			//log.Println("article type:", a.Val)
			this.Type = a.Val
			break
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseArticle(c)
	}

	return
}

func (this *Article) ParseJournal(n *html.Node) (err error) {

	if n.Data == "journal-id" {
		for _, a := range n.Attr {
			if a.Key == "journal-id-type" {
				if a.Val == "nlm-ta" {
					this.Journal.NlmTa = ParseInner(n)
				}
				if a.Val == "iso-abbrev" {
					this.Journal.IsoAbbrev = ParseInner(n)
				}
				if a.Val == "publisher-id" {
					this.Journal.PublisherId = ParseInner(n)
				}
				if a.Val == "hwp" {
					this.Journal.Hwp = ParseInner(n)
				}
				break
			}
		}
	} else if n.Data == "journal-title" {
		title := ParseInner(n)
		this.Journal.Titles = append(this.Journal.Titles, title)
	} else if n.Data == "issn" {
		for _, a := range n.Attr {
			if a.Key == "pub-type" {
				if a.Val == "ppub" {
					this.Journal.Ppub = ParseInner(n)
				}
				if a.Val == "epub" {
					this.Journal.Epub = ParseInner(n)
				}
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseJournal(c)
	}

	return
}

func (this *Article) ParseMeta(n *html.Node) (err error) {

	//log.Println("parse meta:", n.Data)

	if n.Data == "article-id" {
		for _, a := range n.Attr {
			if a.Key == "pub-id-type" {
				if a.Val == "pmid" {
					this.Pmid = ParseInner(n)
				}
				if a.Val == "pmc" {
					this.Pmc = ParseInner(n)
				}
				if a.Val == "doi" {
					this.Doi = ParseInner(n)
				}
				if a.Val == "publisher-id" {
					this.PublisherId = ParseInner(n)
				}
				break
			}
		}
	} else if n.Data == "article-categories" {
		return this.ParseCategories(n)
	} else if n.Data == "title-group" {
		return this.ParseTitles(n)
	} else if n.Data == "contrib-group" {
		return this.ParseContributors(n)
	} else if n.Data == "aff" {
		return this.ParseAff(n)
	} else if n.Data == "author-notes" {
		return this.ParseAuthorNotes(n)
	} else if n.Data == "pub-date" {
		return this.ParsePubDate(n)
	} else if n.Data == "volume" {
		this.Volume = ParseInner(n)
	} else if n.Data == "issue" {
		this.Issue = ParseInner(n)
	} else if n.Data == "fpage" {
		this.Fpage = ParseInner(n)
	} else if n.Data == "lpage" {
		this.Lpage = ParseInner(n)
	} else if n.Data == "history" {
		return this.ParseHistory(n)
	} else if n.Data == "permissions" {
		return this.ParsePermissions(n)
	} else if n.Data == "abstract" {
		this.Abstract, err = this.ParseChildren(n)
	} else if n.Data == "page-count" {
		for _, a := range n.Attr {
			if a.Key == "count" {
				this.PageCount, err = strconv.Atoi(a.Val)
				break
			}
		}
	} else if n.Data == "custom-meta-group" {
		return this.ParseCustomMeta(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseMeta(c)
	}

	//log.Println("META PARSED!!")

	return
}

func (this *Article) ParseBack(n *html.Node) (err error) {

	if n.Data == "ack" {
		this.Ack, err = this.ParseChildren(n)
	} else if n.Data == "ref-list" {
		err = this.ParseRefs(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseBack(c)
	}

	return
}

func (this *Article) ParseCategories(n *html.Node) (err error) {

	if n.Data == "subj-group" {
		for _, a := range n.Attr {
			if a.Key == "subj-group-type" {
				subject := ParseChildInner(n, "subject")
				//log.Println("subject:", subject)
				cat := Category{
					Group:   a.Val,
					Subject: subject,
				}
				this.Categories = append(this.Categories, cat)
				break
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseCategories(c)
	}

	return
}

func (this *Article) ParseTitles(n *html.Node) (err error) {

	if n.Data == "article-title" {
		this.Titles = append(this.Titles, ParseInner(n))
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseTitles(c)
	}
	return
}

func (this *Article) ParseContributors(n *html.Node) (err error) {

	if n.Data == "contrib" {
		this.ParseContributor(n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseContributors(c)
	}
	return
}

func (this *Article) ParseContributor(n *html.Node) (err error) {
	for _, a := range n.Attr {
		if a.Key == "contrib-type" {
			contrib := Contributor{}
			contrib.Type = a.Val
			contrib.Surname = ParseChildInner(n, "surname")
			contrib.GivenNames = ParseChildInner(n, "given-names")
			this.Contributors = append(this.Contributors, contrib)
			break
		}
	}
	return
}

func (this *Article) ParseAff(n *html.Node) (err error) {

	for _, a := range n.Attr {
		if a.Key == "id" {
			this.Aff.Id = a.Val
			break
		}
	}

	this.Aff.Children, err = this.ParseChildren(n)

	return
}

func (this *Article) ParseAuthorNotes(n *html.Node) (err error) {

	if n.Data == "corresp" {
		note := AuthorNote{}
		for _, a := range n.Attr {
			if a.Key == "id" {
				note.CorrespId = a.Val
				break
			}
		}
		note.Children, err = this.ParseChildren(n)
		if err != nil {
			return
		}
		this.AuthorNotes = append(this.AuthorNotes, note)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseAuthorNotes(c)
	}

	return
}

func (this *Article) ParsePubDate(n *html.Node) (err error) {

	for _, a := range n.Attr {
		if a.Key == "pub-type" {
			date := Date{}
			date.Day = ParseChildInner(n, "day")
			date.Month = ParseChildInner(n, "month")
			date.Year = ParseChildInner(n, "year")
			if a.Val == "ppub" {
				this.Ppub = date
			} else if a.Val == "epub" {
				this.Epub = date
			} else if a.Val == "pmc-release" {
				this.PmcRelease = date
			}
		}
	}

	return
}

func (this *Article) ParseHistory(n *html.Node) (err error) {

	if n.Data == "date" {
		date := Date{}
		date.Day = ParseChildInner(n, "day")
		date.Month = ParseChildInner(n, "month")
		date.Year = ParseChildInner(n, "year")
		for _, a := range n.Attr {
			if a.Key == "date-type" {
				if a.Val == "received" {
					this.History.Received = date
				} else if a.Val == "rev-recd" {
					this.History.RevRecd = date
				} else if a.Val == "accepted" {
					this.History.Accepted = date
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseHistory(c)
	}

	return
}

func (this *Article) ParsePermissions(n *html.Node) (err error) {

	if n.Data == "copyright-statement" {
		this.Permissions.CopyrightStatement = ParseInner(n)
	} else if n.Data == "copyright-year" {
		this.Permissions.CopyrightYear = ParseInner(n)
	} else if n.Data == "license" {
		for _, a := range n.Attr {
			if a.Key == "license-type" {
				this.Permissions.Licenses = append(this.Permissions.Licenses, a.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParsePermissions(c)
	}

	return
}

func (this *Article) ParseCustomMeta(n *html.Node) (err error) {

	if n.Data == "custom-meta" {
		meta := Meta{
			Name:  ParseChildInner(n, "meta-name"),
			Value: ParseChildInner(n, "meta-value"),
		}
		this.Metas = append(this.Metas, meta)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseCustomMeta(c)
	}

	return
}

func (this *Article) ParseRefs(n *html.Node) (err error) {

	if n.Data == "title" {
		this.Refs.Title = ParseInner(n)
	}
	if n.Data == "ref" {
		ref := Ref{}
		for _, a := range n.Attr {
			if a.Key == "id" {
				ref.Id = a.Val
			}
		}
		ref.Label = ParseChildInner(n, "label")
		err = ref.ParseGroups(n)
		err = ref.ParseElement(n)

		this.Refs.List = append(this.Refs.List, ref)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseRefs(c)
	}

	return
}

func (this *Ref) ParseGroups(n *html.Node) (err error) {

	if n.Data == "person-group" {
		group := Group{}
		for _, a := range n.Attr {
			if a.Key == "person-group-type" {
				group.Type = a.Val
			}
		}
		err = group.ParseNames(n)
		if err != nil {
			return
		}
		this.Groups = append(this.Groups, group)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		this.ParseGroups(c)
	}

	return
}

func (this *Group) ParseNames(n *html.Node) (err error) {

	if n.Data == "name" {
		name := map[string]string{}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			inner := ParseInner(c)
			if inner != "" {
				name[c.Data] = inner
			}
		}
		//log.Println("name:", name)
		this.Names = append(this.Names, name)
		//log.Println("names:", this.Names)
		return
	}

	//log.Println("names:", names)

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err = this.ParseNames(c)
	}

	return
}

func (this *Ref) ParseElement(n *html.Node) (err error) {

	if n.Data == "element-citation" {
		//log.Println("element:", n.Data)
		this.Title = ParseChildInner(n, "article-title")
		this.Source = ParseChildInner(n, "source")
		this.Year = ParseChildInner(n, "year")
		this.Volume = ParseChildInner(n, "volume")
		this.Fpage = ParseChildInner(n, "fpage")
		this.Lpage = ParseChildInner(n, "lpage")
		this.Pmid = ParseChildInner(n, "pub-id")
		return
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		err = this.ParseElement(c)
	}

	return
}

func ParseChildInner(n *html.Node, child string) (val string) {
	if n.Data == child {
		val = ParseInner(n)
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		v := ParseChildInner(c, child)
		if v != "" {
			val = v
		}
	}
	return
}

func ParseInner(n *html.Node) (inner string) {
	c := n.FirstChild
	//log.Println("parse inner", n.Data, n.Attr, c.Data)
	if c == nil {
		return ""
	}
	inner = c.Data
	return
}

func (this *Article) ParseChildren(n *html.Node) (children []Node, err error) {

	//read in each child node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		child := Node{}
		if c.Type == html.TextNode {
			child.Type = "text"
			child.Body = c.Data
		} else {
			child.Type = "tag"
			child.Tag = c.Data
			child.Props = map[string]string{}
			for _, a := range c.Attr {
				//log.Println("key:", a.Key, child.Props)
				child.Props[a.Key] = a.Val
			}
			child.Children, err = this.ParseChildren(c)
			if err != nil {
				return
			}
			err = child.GetSentences()
			if err != nil {
				return
			}

		}
		children = append(children, child)
		this.ParseContributors(c)
	}

	/*
		var buff bytes.Buffer
		err = html.Render(&buff, n)
		if err != nil {
			return
		}

		htmlString = buff.String()
	*/
	return
}

func (this *Node) GetSentences() (err error) {
	//loop through children and get sentences start and end

	body := ""
	for _, child := range this.Children {
		if child.Type == "text" {
			body += child.Body
		}
	}

	//split based on sentences
	i := 0
	curI := 0
	for i != -1 {
		i = strings.Index(body, ". ")
		sent := Sentence{}
		if i == -1 { //none found
			sent.Start = curI
			sent.End = curI + len(body)
		} else {
			sent.Start = curI
			sent.End = curI + i + 1
		}
		this.Sentences = append(this.Sentences, sent)
		if i == -1 {
			break
		}
		body = body[i+1:]
		curI += i + 1
	}

	return
}
