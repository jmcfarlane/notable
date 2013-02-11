/**
 * @fileoverview Describes a bootstrap tab
 */
define([
  'text!templates/noteDetail.html',
  'text!templates/tab.html',
  'backbone',
  'underscore',
  'lib/jquery.hotkeys',
  'codemirror',
  'lib/codemirror/mode/rst/rst'
],
function(noteDetailTemplate, tabTemplate) {
  return Backbone.View.extend({

    initialize: function() {
      this._editor = null;
    },

    render: function(collection) {
      var note = this.model.toJSON(),
        tabs = this.options.tabs,
        tabContent = this.options.tabContent;

      // No tabs are active
      tabs.find('.active').removeClass('active');

      // Add new tab content
      tabContent.append(_.template(noteDetailTemplate, {
        note: note
      }));

      // Use codemirror for the content
      this.el = $(tabContent).find('.editor').last().parent();
      this.$el = $(this.el);
      this._editor = CodeMirror(_.first(this.$el.find('.editor')), {
        mode: 'rst',
        value: note.content,
        extraKeys: {
          'Ctrl-S': _.bind(this.save, this)
        }
      });
      this.$('input').bind('keydown', 'ctrl+s', _.bind(this.save, this));

      // Somehow this seems to result in the right tab saving, but
      // seemms like a defect waiting to happen.
      $(document).bind('keydown', 'ctrl+s', _.bind(this.save, this));

      // Add a new tab
      tabs.append(_.template(tabTemplate, {
        note: note
      }));

      // Set event handlers, and show the tab
      this.getTab().on('shown', _.bind(this.shown, this))
        .tab('show')
        .find('button').on('click', _.bind(this.close, this));

      return this;
    },

    show: function() {
      this.getTab().tab('show');
    },

    shown: function() {
      if (this.model.get('subject')) {
        this._editor.focus();
      } else {
        this.$('.subject input').focus();
      }
    },

    saved: function() {
      $('.saved').fadeIn();
      setTimeout(function() {
        $('.saved').fadeOut();
      }, 2000);
    },

    save: function() {
      if (!this.$el.is(":visible")) {
        return;
      }
      this.model.save({
        content: this._editor.getValue(),
        password: this.$el.find('.password input').val(),
        subject: this.$el.find('.subject input').val(),
        tags: this.$el.find('.tags input').val()
      }, {
        success: this.saved
      });
      return false;
    },

    getTab: function() {
      return this.options.tabs.find(this.selector());
    },

    selector: function() {
      return 'a[href=#'+ this.model.get('uid') +']';
    },

    close: function() {
      // Show the previous tab and then remove this one
      this.options.tabs.find(this.selector())
        .parent()
        .prev()
        .find('a')
        .tab('show')
        .parent()
        .next()
        .remove();
      // Remove the tab content (after the ^ animation is finished) (ick)
      setTimeout(_.bind(function() {
        $('#' + this.model.get('uid')).detach().remove();
      }, this), 300);

      // Let parents deal with me
      this.trigger('destroy');
    }

  });
});
