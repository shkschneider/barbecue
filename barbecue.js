// globals

$.store = {
    get: function (key) {
        return window.sessionStorage.getItem(key);
    },
    set: function (key, value) {
        if (value) {
            window.sessionStorage.setItem(key, value);
        } else {
            window.sessionStorage.removeItem(key);
        }
    },
    clear: function () {
        window.sessionStorage.clear();
    },
    keys: function () {
        return Object.keys(window.sessionStorage);
    },
};

// functions

$.id = function () {
    return Math.floor(Math.random() * Math.floor(Math.random() * Date.now())).toString(16);
};

$.slug = function (s) {
    return s.replace(/[^a-z0-9]/gi, '-')
};

$.sanitize = function (s) {
    return s.replaceAll(/<\/?[^>]+(>|$)/gi, '');
};

// //debug
// if ($.store.keys().length == 0) {
//     $.store.set('1', JSON.stringify({ 'id': '1',
//         'title': 'Avant-Projet', 'description': 'Avant!' }));
//     $.store.set('11', JSON.stringify({ 'id': '11',
//         'parent': '1', 'title': 'Etude', 'description': 'etude de...' }));
//     $.store.set('12', JSON.stringify({ 'id': '12',
//         'parent': '1', 'title': 'Plans', 'description': 'plans de...' }));
//     $.store.set('121', JSON.stringify({ 'id': '121',
//         'parent': '12', 'title': 'masse', 'progress': 100 }));
//     $.store.set('122', JSON.stringify({ 'id': '122',
//         'parent': '12', 'title': 'technique' }));
//     $.store.set('13', JSON.stringify({ 'id': '13',
//         'parent': '1', 'title': 'Arbitrage', 'description': 'arbitrage de...' }));
//     $.store.set('14', JSON.stringify({ 'id': '14',
//         'parent': '1', 'title': 'Permis', 'description': 'permis de...' }));
//     $.store.set('2', JSON.stringify({ 'id': '2',
//         'title': 'Chantier' }));
// }
// //debug

$.barbecue = {
    getById: function (id) {
        let o = JSON.parse($.store.get(id));
        if (o) {
            console.assert(o.id, 'id null');
            console.assert(o.title, 'title null');
            o.progress = Number(o.progress) || 0;
        }
        return o;
    },
    getAll: function () {
        let a = new Array();
        for (let k of $.store.keys()) {
            let o = $.barbecue.getById(k);
            a.push($.barbecue.sauce(o));
        }
        return a;
    },
    getRoots: function () {
        let a = new Array();
        for (let k of $.store.keys()) {
            let o = $.barbecue.getById(k);
            if (!o.parent) {
                a.push($.barbecue.sauce(o));
            }
        }
        return a;
    },
    parentsOf: function (c) {
        let a = new Array();
        while (c.parent) {
            let p = $.barbecue.getById(c.parent);
            if (!p) {
                break;
            }
            a.push(p);
            c = p;
        }
        return a.reverse();
    },
    childrenOf: function (p) {
        let a = new Array();
        for (let k of $.store.keys()) {
            let c = $.barbecue.getById(k);
            if (c.parent == p.id) {
                a.push($.barbecue.sauce(c));
            }
        }
        return a;
    },
    sauce: function (o) {
        o.title = o.title || '?';
        o.description = o.description || '';
        let a = new Array();
        for (let k of $.store.keys()) {
            let c = $.barbecue.getById(k);
            if (c.parent == o.id) {
                a.push(Number($.barbecue.sauce(c).progress) || 0);
            }
        }
        if (a.length) {
            o.progress = a.reduce((sum, c) => sum + c || 0, 0) / a.length;
            o.progress = Math.min(Math.max(o.progress || 0, 0), 100).toFixed(0);
        } else {
            o.progress = Number(o.progress || 0);
        }
        o.progress = o.progress || 0;
        return o;
    },
    import: function () {
        let file = document.getElementById('file').files[0];
        if (file.length <= 0) { return ; }
        let reader = new FileReader();
        reader.onload = function (e) {
            let json = JSON.parse(e.target.result);
            $.each(json, function (key) {
                let data = json[key];
                $.store.set(data.id, JSON.stringify(data));
            });
            window.location.href = 'index.html';
        }
        reader.readAsText(file);
    },
    export: function () {
        let json = Object.entries(window.sessionStorage).reduce(
            (obj, [k, v]) =>  ({...obj, [k]: JSON.parse(v)}),
            {}
        )
        $('.header #actions a#export').attr("href", "data:text/json;charset=utf-8," + encodeURIComponent(JSON.stringify(json, undefined, 2)));
        $('.header #actions a#export').attr("download", "barbecue.json");
        document.querySelector('a#export').click();
    },
    clear: function () {
        $.store.clear();
        window.location.href = 'index.html';
    },
    add: function () {
        let parent = ($.o) ? ($.o.id || null) : null;
        let title = $('article footer #title').val() || '?';
        let description = $('article footer #description').val() || '';
        const id = $.id();
        const slug = $.slug(title);
        $.store.set(id, JSON.stringify({
            'id': encodeURIComponent(id),
            'slug': encodeURIComponent(slug),
            'parent': parent,
            'title': title,
        }));
        if ($.o) {
            window.location.href = `index.html?${$.o.id}=${$.o.slug}`;
        } else {
            window.location.href = `index.html?${id}=${slug}`;
        }
    },
    title: function () {
        if ($('article header input#edit-title').hasClass('gone')) {
            $('article header nav#navigation').hide();
            $('article header input#edit-title').prop('value', $.o.title);
            $('article header input#edit-title').removeClass('gone');
            $('article header a#edit').addClass('gone');
            $('article header a#save').removeClass('gone');
            $('article header a#delete').removeClass('gone');
            $('article header a#complete').addClass('gone');
            $('article header progress#progress').addClass('gone');
            $('article header input#edit-progress').prop('value', $.o.progress);
            $('article header input#edit-progress').removeClass('gone');
            // }
        } else {
            $.o.title = $.sanitize($('article header input#edit-title').val());
            $.o.slug = encodeURIComponent($.slug($.o.title));
            $.o.progress = Number($('article header input#progress').val() || $.o.progress).toFixed(0);
            $.store.set($.o.id, JSON.stringify($.o));
            $('article header nav#navigation ul li:last-child').html(`${$.o.title} (${$.o.progress}%)`);
            $('article header progress#progress').prop('value', $.o.progress);
            $('article header nav#navigation').show();
            $('article header input#edit-title').addClass('gone');
            $('article header a#edit').removeClass('gone');
            $('article header a#save').addClass('gone');
            $('article header a#delete').addClass('gone');
            if ($.barbecue.childrenOf($.o).length == 0) {
                $('article header a#complete').removeClass('gone');
            }
            $('article header progress#progress').prop('value', $.o.progress);
            $('article header progress#progress').removeClass('gone');
            $('article header input#edit-progress').addClass('gone');
            window.location.href = `index.html?${$.o.id}=${$.o.slug}`;
        }
    },
    description: function () {
        if ($('#content #description textarea').length == 0) {
            let lines = $.o.description.split(/\r\n|\r|\n/).length + 1;
            $('#content #description p:first-child').html(`<textarea rows="${lines}" placeholder="Description...">${$.o.description}</textarea>`);
            $('#content #description a#edit-description').html('Save');
        } else {
            $.o.description = $.sanitize($('#content #description p:first-child textarea').val() || '');
            $.store.set($.o.id, JSON.stringify($.o));
            $('#content #description p:first-child').html(marked.parse($.o.description));
            $('#content #description a#edit-description').html('Edit');
        }
    },
    complete: function () {
        $.o.progress = 100;
        $.store.set($.o.id, JSON.stringify($.o));
        window.location.href = `index.html?${$.o.parent}=`;
    },
    delete: function () {
        $.store.set($.o.id, null);
        if ($.o.parent) {
            window.location.href = `index.html?${$.o.parent}=`;
        } else {
            window.location.href = 'index.html';
        }
    },
};
// everything coming out of barbecue should be sauced
let query = new URLSearchParams(location.search);
if (query.size) {
    for (let q of query) {
        $.o = $.barbecue.getById(q[0]);
    }
    if ($.o) {
        $.o = $.barbecue.sauce($.o);
    } else {
        alert(`Object not found: ${query}`);
        window.location.href = 'index.html';
    }
}

// header

$(function () {
    if ($.o) {
        for (let p of $.barbecue.parentsOf($.o)) {
            $('article header nav#navigation ul').append(`<li><a href="index.html?${p.id}=${p.slug}" class="contrast">${p.title}</a></li>`);
        }
        $('article header nav#navigation ul').append(`<li>${$.o.title} (${$.o.progress}%)</li>`);
        $('article header a#edit').removeClass('gone');
        if ($.o.progress < 100 && $.barbecue.childrenOf($.o).length == 0) {
            $('article header a#complete').removeClass('gone');
        }
        $('article header #progress').prop('value', $.o.progress);
    } else {
        if ($.barbecue.getAll().length == 0) {
            $('.header li#import').toggleClass('gone');
            $('.header li#export').remove();
            $('.header li#clear').remove();
            let quotes = [
                "Life is short. Do stuff that matters. — Siqi Chen",
                "It always seems impossible until it’s done. — Nelson Mandela",
                "Productivity is never an accident. — Paul J. Meyer",
                "Amateurs sit and wait for inspiration, the rest of us just get up and go to work. — Stephen King",
                "If you spend too much time thinking about a thing, you’ll never get it done. — Bruce Lee",
                "Your mind is for having ideas, not holding them. — David Allen",
                "You may delay, but time will not. — Benjamin Franklin",
                "While one person hesitates because he feels inferior, the other is busy making mistakes and becoming superior. — Henry Link",
                "The only way around is through. — Robert Frost",
                "The secret of getting things done is to act! — Dante Alighieri",
                "Procrastination is opportunity’s assassin. — Victor Kiam",
                "Only put off until tomorrow what you are willing to die having left undone. — Pablo Picasso",
                "You can do anything, but not everything. — David Allen",
                "The great secret about goals and visions is not the future they describe but the change in the present they engender.",
                "Use your mind to think about things, rather than think of them.",
                "Your ability to generate power is directly proportional to your ability to relax.",
                "Pick battles big enough to matter, small enough to win. — Jonathan Kozol",
                "It does not take much strength to do things, but it requires a great deal of strength to decide what to do. — Elbert Hubbard"
            ]
            $('article section#children').html('<blockquote>' + quotes[Math.floor(Math.random() * quotes.length)] + '</blockquote>');
        } else {
            $('.header li#import').remove();
            $('.header li#export').toggleClass('gone');
            $('.header li#clear').toggleClass('gone');
        }
        $('article header #progress').hide();
    }
});

// content

$(function () {
    if ($.o) {
        if ($.o.description) {
            $('#content #description p:first-child').html(marked.parse($.o.description));
        }
        let children = $.barbecue.childrenOf($.o);
        if (children.length) {
            let n = children.filter((c) => c.progress == 100).length;
            let N = children.length;
            $('#content #children').append(`<p class="count">${n}/${N}</p>`);
            let ul = $('<ul class="children"></ul>');
            for (let c of children) {
                let li = $(`<li class="child progress-${c.progress}"></li>`);
                    let div = $('<div class="grid"></div>');
                        $(div).append(`<a href="index.html?${c.id}=${c.slug}" class="contrast">${c.title}</a>`);
                        $(div).append(`<span style="text-align:right">${c.progress}%</span>`);
                    $(li).append(div);
                    $(li).append('<br />');
                    $(li).append(`<progress value="${c.progress}" max="100"></progress>`);
                $(ul).append(li);
            }
            $('#content #children').append(ul);
        } else {
            $('#content #children').hide();
        }
    } else {
        $('#content a#edit').hide();
        $('#content a#progress').hide();
        $('#content #description').hide();
        let ul = $('<ul class="children"></ul>');
        for (let o of $.barbecue.getRoots()) {
            let li = $(`<li class="child progress-${o.progress}"></li>`);
                let div = $('<div class="grid"></div>');
                    $(div).append(`<a href="index.html?${o.id}=${o.slug}" class="contrast">${o.title}</a>`);
                    $(div).append(`<span style="text-align:right">${o.progress}%</span>`);
                $(li).append(div);
                $(li).append('<br />');
                $(li).append(`<progress value="${o.progress}" max="100"></progress>`);
            $(ul).append(li);
        }
        $('#content #children').append(ul);
    }
    $('#content').removeAttr('aria-busy');
});

// footer

if ($.o) {
    $('footer form select').append(`<option value="${$.o.id}" selected>${$.o.title}</option>`);
}

// EOF
