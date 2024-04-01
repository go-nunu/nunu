import{_ as s,c as i,o as a,a4 as n}from"./chunks/framework.DZjeu1b3.js";const F=JSON.parse('{"title":"Model","description":"","frontmatter":{},"headers":[],"relativePath":"model.md","filePath":"model.md","lastUpdated":1711880770000}'),t={name:"model.md"},l=n(`<h1 id="model" tabindex="-1">Model <a class="header-anchor" href="#model" aria-label="Permalink to &quot;Model&quot;">​</a></h1><p>在一个 Go 语言项目中使用 GORM (GORM) 时，通常会将数据库表映射到 Go 的结构体。这些结构体定义了数据模型，用于表示数据库中的表和表之间的关系。在 GORM 项目中，通常会将这些结构体定义在一个名为 model 的目录中。</p><h2 id="核心用途" tabindex="-1">核心用途 <a class="header-anchor" href="#核心用途" aria-label="Permalink to &quot;核心用途&quot;">​</a></h2><p>这个目录的作用主要是组织和管理与数据库模型相关的代码。将所有的模型文件放置在一个目录下可以更好地组织代码，并使代码结构清晰易于维护。通常，每个模型都会对应一个结构体，该结构体的字段与数据库表的列相对应。</p><p>除了结构体定义之外，model 目录中的文件可能还包括与模型相关的其他代码，如模型的方法、模型之间的关系定义等。</p><p>下面是一个简单的 User 结构体定义的例子：</p><div class="language-go vp-adaptive-theme"><button title="Copy Code" class="copy"></button><span class="lang">go</span><pre class="shiki shiki-themes github-light github-dark vp-code"><code><span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">type</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;"> User</span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;"> struct</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> {</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	Id        </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">uint</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;">   \`gorm:&quot;primarykey&quot;\`</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	UserId    </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">string</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> \`gorm:&quot;unique;not null&quot;\`</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	Nickname  </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">string</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> \`gorm:&quot;not null&quot;\`</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	Password  </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">string</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> \`gorm:&quot;not null&quot;\`</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	Email     </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">string</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> \`gorm:&quot;not null&quot;\`</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	CreatedAt </span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">time</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Time</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	UpdatedAt </span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">time</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">Time</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">	DeletedAt </span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">gorm</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">.</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">DeletedAt</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> \`gorm:&quot;index&quot;\`</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span>
<span class="line"></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">func</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> (</span><span style="--shiki-light:#E36209;--shiki-dark:#FFAB70;">u </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">*</span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">User</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">) </span><span style="--shiki-light:#6F42C1;--shiki-dark:#B392F0;">TableName</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">() </span><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">string</span><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;"> {</span></span>
<span class="line"><span style="--shiki-light:#D73A49;--shiki-dark:#F97583;">	return</span><span style="--shiki-light:#032F62;--shiki-dark:#9ECBFF;"> &quot;users&quot;</span></span>
<span class="line"><span style="--shiki-light:#24292E;--shiki-dark:#E1E4E8;">}</span></span></code></pre></div>`,7),h=[l];function e(p,k,r,d,E,o){return a(),i("div",null,h)}const y=s(t,[["render",e]]);export{F as __pageData,y as default};